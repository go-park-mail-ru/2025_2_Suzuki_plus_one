package dto

//go:generate easyjson -all $GOFILE
type GetUserMeInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
}


//go:generate easyjson -all $GOFILE
type GetUserMeOutput struct {
	ID                 uint     `json:"id"`
	Username           string   `json:"username"`
	Email              string   `json:"email"`
	DateOfBirth        JSONDate `json:"date_of_birth"`
	PhoneNumber        string   `json:"phone_number"`
	AvatarURL          string   `json:"avatar_url"`
	SubscriptionStatus string   `json:"subscription_status"`
}
