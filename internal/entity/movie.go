package entity

import (
	"fmt"
)

type Movie struct {
	ID      string   `json:"id"`
	Title   string   `json:"title" validate:"required"`
	Year    int      `json:"year"`
	Genres  []string `json:"genres"`
	Preview string   `json:"preview"`
}

func NewMovie(id, title string, year int, genres []string, preview string) (Movie, error) {
	m := Movie{
		ID:      id,
		Title:   title,
		Year:    year,
		Genres:  genres,
		Preview: preview,
	}
	if err := validate.Struct(m); err != nil {
		return Movie{}, fmt.Errorf("invalid movie entity validation: %w", err)
	}
	return m, nil
}

func NewTestMovie(id string) Movie {
	return Movie{
		ID:      id,
		Title:   "Test Movie",
		Year:    2023,
		Genres:  []string{"Drama", "Action"},
		Preview: "http://example.com/preview.jpg",
	}
}
