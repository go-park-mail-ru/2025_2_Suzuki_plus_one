package http

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

const PRICE = "1.00" // TODO: remove hardcoded price

type PostPaymentNewUsecase struct {
	logger      logger.Logger
	paymentRepo PaymentRepository
}

func NewPostPaymentNewUsecase(
	logger logger.Logger,
	paymentRepo PaymentRepository,
) *PostPaymentNewUsecase {
	if logger == nil {
		panic("NewPostPaymentNewUsecase: logger is nil")
	}
	if paymentRepo == nil {
		panic("NewPostPaymentNewUsecase: paymentRepo is nil")
	}
	return &PostPaymentNewUsecase{
		logger:      logger,
		paymentRepo: paymentRepo,
	}
}

func (uc *PostPaymentNewUsecase) Execute(
	ctx context.Context,
	input dto.PostPaymentNewInput,
) (
	dto.PostPaymentNewOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("PostPaymentNewUsecase called")

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInParamsInvalid,
			err.Error(),
		)
		log.Error("Invalid sign in input parameters", log.ToError(err))
		return dto.PostPaymentNewOutput{}, &derr
	}

	// Get access token info
	userID, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_payment_new",
			entity.ErrPostPaymentNewAccessTokenInvalid,
			"access token is invalid: "+err.Error(),
		)
		log.Error("Access token is invalid", log.ToError(err))
		return dto.PostPaymentNewOutput{}, &derr
	}

	// Create payment in payment repository
	paymentID, err := uc.paymentRepo.CreatePayment(
		ctx,
		userID,
		PRICE,
		"Payment for user ID "+strconv.Itoa(int(userID)),
	)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_payment_new",
			entity.ErrPostPaymentNewAccessTokenInvalid,
			"paymentRepo.CreatePayment failed: "+err.Error(),
		)
		log.Error("Failed to create payment", log.ToError(err))
		return dto.PostPaymentNewOutput{}, &derr
	}

	log.Debug("Payment created successfully", log.ToString("paymentID", paymentID))

	// Return output DTO
	return dto.PostPaymentNewOutput{
		PaymentID: paymentID,
	}, nil
}
