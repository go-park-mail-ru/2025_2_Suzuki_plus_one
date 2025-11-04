package entity

type User struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	AssetImageID uint
}
