package yookassa

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

func (yk *Yookassa) CreatePayment(ctx context.Context, userID uint, amount string, description string) (*yoopayment.Payment, error) {
	log := logger.LoggerWithKey(yk.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CreatePayment called",
		log.ToString("amount", amount),
		log.ToString("description", description),
	)

	handler := yookassa.NewPaymentHandler(yk.Client)
	// Создаем платеж
	payment, err := handler.CreatePayment(&yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    amount,
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType("bank_card"),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: yk.redirectURL,
		},
		Description:       description,
		SavePaymentMethod: false,
		Capture:           false,
		// Metadata:           map[string]string{"user_id": strconv.Itoa(int(userID))},
		MerchantCustomerID: strconv.Itoa(int(userID)),
	})

	if err != nil {
		log.Error("CreatePayment: failed to create payment", log.ToError(err))
		return nil, err
	}

	log.Info("Payment created successfully",
		log.ToString("payment.ID", payment.ID),
		log.ToAny("payment.Test", payment.Test),
		log.ToAny("payment", payment),
	)

	yk.Handlers[payment.ID] = handler
	return payment, nil
}

func (yk *Yookassa) CapturePayment(ctx context.Context, payment *yoopayment.Payment) (*yoopayment.Payment, error) {
	log := logger.LoggerWithKey(yk.logger, ctx, common.ContextKeyRequestID)
	log.Debug("CapturePayment called",
		log.ToString("payment.ID", payment.ID),
	)

	handler, ok := yk.Handlers[payment.ID]
	if !ok {
		log.Error("CapturePayment: payment handler not found for payment ID", log.ToString("payment.ID", payment.ID))
		return nil, errors.New("CapturePayment: payment handler not found")
	}

	capturedPayment, err := handler.CapturePayment(payment)

	if err != nil {
		log.Error("CapturePayment: failed to capture payment", log.ToError(err))
		return nil, err
	}

	log.Info("Payment captured successfully",
		log.ToString("capturedPayment.ID", capturedPayment.ID),
		log.ToAny("capturedPayment.Status", capturedPayment.Status),
	)
	delete(yk.Handlers, payment.ID)
	
	return capturedPayment, nil
}
