package http

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
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
	payment, err := uc.paymentRepo.CreatePayment(
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

	log.Debug("Payment created successfully", log.ToString("paymentID", payment.ID))

	// Extract confirmation URL from different possible types
	var confirmationURLStr string
	switch c := payment.Confirmation.(type) {
	case *yoopayment.Redirect:
		confirmationURLStr = c.ConfirmationURL
	case yoopayment.Redirect:
		confirmationURLStr = c.ConfirmationURL
	case map[string]any:
		// try common JSON keys
		if v, ok := c["confirmation_url"].(string); ok && v != "" {
			confirmationURLStr = v
		} else if v, ok := c["confirmationUrl"].(string); ok && v != "" {
			confirmationURLStr = v
		}
	default:
		// Fallback: try marshaling into Redirect struct
		var redirect yoopayment.Redirect
		if b, err := json.Marshal(payment.Confirmation); err == nil {
			if err2 := json.Unmarshal(b, &redirect); err2 == nil && redirect.ConfirmationURL != "" {
				confirmationURLStr = redirect.ConfirmationURL
			}
		}
	}

	if confirmationURLStr == "" {
		derr := dto.NewError(
			"usecase/post_payment_new",
			entity.ErrPostPaymentNewConfirmationURLInvalid,
			"payment.Confirmation does not contain a confirmation URL",
		)
		log.Error("payment.Confirmation does not contain a confirmation URL", log.ToAny("payment.Confirmation", payment.Confirmation))
		return dto.PostPaymentNewOutput{}, &derr
	}

	// Return output DTO
	return dto.PostPaymentNewOutput{
		RedirectURL: confirmationURLStr,
		PaymentID:       payment.ID,
	}, nil
}
