package entity

import "time"

type RefreshToken struct {
	ID        uint
	Token     string
	UserID    uint
	ExpiresAt time.Time
}
