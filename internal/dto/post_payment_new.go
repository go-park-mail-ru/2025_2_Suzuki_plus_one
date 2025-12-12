package dto

type PostPaymentNewInput struct {
	AccessToken string `json:"access_token"`
}

type PostPaymentNewOutput struct {
	PaymentID       string `json:"payment_id"`
	RedirectURL string `json:"redirect_url"`
}