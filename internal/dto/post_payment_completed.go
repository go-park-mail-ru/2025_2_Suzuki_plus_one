package dto

import (
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
	yoowebhook "github.com/rvinnie/yookassa-sdk-go/yookassa/webhook"
)

type PostPaymentCompletedInput struct {
	Webhook yoowebhook.WebhookEvent[yoopayment.Payment] `json:"yoopayment"`
}

type PostPaymentCompletedOutput struct {
	Status string `json:"status"`
}

// //go:generate easyjson -all post_payment_completed.go
