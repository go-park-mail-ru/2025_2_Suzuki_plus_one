package entity

import "time"

type Actor struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birth_date"`
	Bio       string    `json:"bio"`
}
