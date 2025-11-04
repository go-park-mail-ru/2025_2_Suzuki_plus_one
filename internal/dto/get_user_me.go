package dto

import "time"

type GetUserMeInput struct {
	AccessToken string `json:"access_token"` // Access token from Authorization header
}

type GetUserMeOutput struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
	PhoneNumber string    `json:"phone_number"`
	AvatarURL   string    `json:"avatar_url"`
}
