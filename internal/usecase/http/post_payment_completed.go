package http

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	yoowebhook "github.com/rvinnie/yookassa-sdk-go/yookassa/webhook"
)

type PostPaymentCompletedUsecase struct {
	logger   logger.Logger
	userRepo UserRepository
	payRepo  PaymentRepository
}

func NewPostPaymentCompletedUsecase(
	logger logger.Logger,
	userRepo UserRepository,
	payRepo PaymentRepository,
) *PostPaymentCompletedUsecase {
	if logger == nil {
		panic("NewPostPaymentCompletedUsecase: logger is nil")
	}
	if userRepo == nil {
		panic("NewPostPaymentCompletedUsecase: userRepo is nil")
	}
	if payRepo == nil {
		panic("NewPostPaymentCompletedUsecase: payRepo is nil")
	}
	return &PostPaymentCompletedUsecase{
		logger:   logger,
		userRepo: userRepo,
		payRepo:  payRepo,
	}
}

// This usecase is hit when the payment gateway notifies us that the payment is completed.
func (uc *PostPaymentCompletedUsecase) Execute(
	ctx context.Context,
	input dto.PostPaymentCompletedInput,
) (
	dto.PostPaymentCompletedOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("PostPaymentCompletedUsecase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInParamsInvalid,
			err.Error(),
		)
		log.Error("Invalid sign in input parameters", log.ToError(err))
		return dto.PostPaymentCompletedOutput{}, &derr
	}

	userID, err := strconv.Atoi(input.Webhook.Object.MerchantCustomerID)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_payment_completed",
			entity.ErrPostPaymentCompletedInvalidParams,
			"Invalid MerchantCustomerID in webhook",
		)
		log.Error("Invalid MerchantCustomerID in webhook", log.ToError(err))
		return dto.PostPaymentCompletedOutput{}, &derr
	}

	switch input.Webhook.Event {
	case yoowebhook.EventPaymentWaitingForCapture:
		log.Info("Processing payment waiting for capture event",
			log.ToString("paymentID", input.Webhook.Object.ID),
		)
		_, err := uc.payRepo.CapturePayment(ctx, &input.Webhook.Object)
		if err != nil {
			derr := dto.NewError(
				"usecase/post_payment_completed",
				entity.ErrPostPaymentCompletedInvalidParams,
				"Failed to capture payment: "+err.Error(),
			)
			log.Error("Failed to capture payment", log.ToError(err))
			return dto.PostPaymentCompletedOutput{}, &derr
		}
		log.Debug("Payment captured successfully",
			log.ToString("paymentID", input.Webhook.Object.ID),
		)
		// You might also want to notify the user that their payment was successful.
		if uerr := uc.userRepo.UpdateUserSubscriptionStatus(ctx, uint(userID), "pending"); uerr != nil {
			log.Error("Failed to set user subscription to pending", log.ToError(uerr))
		}
		log.Debug("User subscription pending status set",
			log.ToInt("userID", userID),
		)
	case yoowebhook.EventPaymentSucceeded:
		log.Info("Processing payment succeeded event",
			log.ToString("paymentID", input.Webhook.Object.ID),
		)
		if uerr := uc.userRepo.UpdateUserSubscriptionStatus(ctx, uint(userID), "active"); uerr != nil {
			log.Error("Failed to set user subscription to active", log.ToError(uerr))
		}
		log.Debug("User subscription active status set",
			log.ToInt("userID", userID),
		)
	case yoowebhook.EventPaymentCanceled:
		log.Warn("Processing payment canceled event",
			log.ToString("paymentID", input.Webhook.Object.ID),
		)
		if uerr := uc.userRepo.UpdateUserSubscriptionStatus(ctx, uint(userID), "inactive"); uerr != nil {
			log.Error("Failed to set user subscription to inactive", log.ToError(uerr))
		}
		log.Debug("User subscription inactive status set",
			log.ToInt("userID", userID),
		)
	default:
		log.Error("Unhandled webhook event type",
			log.ToString("event", string(input.Webhook.Event)),
		)
	}

	// Return output DTO
	return dto.PostPaymentCompletedOutput{}, nil
}
