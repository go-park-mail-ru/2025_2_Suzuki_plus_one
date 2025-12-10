package yookassa

import (
	"context"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

func (yk *Yookassa) CreatePayment(ctx context.Context, userID uint, amount string, description string) (string, error) {
	log := logger.LoggerWithKey(yk.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CreatePayment called",
		log.ToString("amount", amount),
		log.ToString("description", description),
	)

	// Создаем платеж
	payment, err := yk.Handler.CreatePayment(&yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    amount,
			Currency: "RUB",
		},
		// PaymentMethod: yoopayment.PaymentMethodType("bank_card"),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: yk.redirectURL,
		},
		Description:        "Test payment",
		MerchantCustomerID: strconv.Itoa(int(userID)),
	})

	if err != nil {
		log.Error("CreatePayment: failed to create payment", log.ToError(err))
		return "", err
	}

	log.Info("Payment created successfully",
		log.ToString("payment.ID", payment.ID),
		log.ToAny("payment.Test", payment.Test),
	)

	return payment.ID, nil
}

func (yk *Yookassa) CapturePayment(ctx context.Context, payment *yoopayment.Payment) (*yoopayment.Payment, error) {
	log := logger.LoggerWithKey(yk.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CapturePayment called",
		log.ToString("payment.ID", payment.ID),
	)

	capturedPayment, err := yk.Handler.CapturePayment(payment)

	if err != nil {
		log.Error("CapturePayment: failed to capture payment", log.ToError(err))
		return nil, err
	}

	log.Info("Payment captured successfully",
		log.ToString("capturedPayment.ID", capturedPayment.ID),
		log.ToAny("capturedPayment.Status", capturedPayment.Status),
	)

	return capturedPayment, nil
}
