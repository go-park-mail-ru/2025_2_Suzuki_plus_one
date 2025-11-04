package entity

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	DateOfBirth  string `json:"date_of_birth"`
	PhoneNumber  string `json:"phone_number"`
	AssetImageID uint   `json:"asset_image_id"`
}
