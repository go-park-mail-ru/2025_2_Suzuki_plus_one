package dto

type PostPaymentNewInput struct {
	AccessToken string `json:"access_token"`
}

type PostPaymentNewOutput struct {
	PaymentID       string `json:"payment_id"`
	ConfirmationURL string `json:"confirmation_url"`
}