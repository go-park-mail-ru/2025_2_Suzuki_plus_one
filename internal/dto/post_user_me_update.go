package dto

type PostUserMeUpdateInput struct {
	AccessToken string   `json:"access_token" required:"true"` // Access token from Authorization header
	Username    string   `json:"username" example:"johndoe" required:"true"`
	Email       string   `json:"email" format:"email" example:"johndoe@example.com" required:"true"`
	DateOfBirth JSONDate `json:"date_of_birth" format:"date" example:"1990-01-01" required:"true"`
	PhoneNumber string   `json:"phone_number" example:"+1234567890" required:"true"`
}

type PostUserMeUpdateOutput struct {
	GetUserMeOutput
}
