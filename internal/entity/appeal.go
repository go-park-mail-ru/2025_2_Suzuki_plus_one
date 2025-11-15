package entity

import "time"

type Appeal struct {
	ID        uint      `json:"appeal_id"`
	UserID    uint      `json:"user_id"`
	Tag       string    `json:"tag"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type AppealMessage struct {
	ID         uint      `json:"appeal_message_id"`
	IsResponse bool      `json:"is_response"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"-"`
}
