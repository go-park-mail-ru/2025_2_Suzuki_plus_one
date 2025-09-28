package db

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
)

func TestFindUserByEmail(t *testing.T) {
	db := NewMockDB()

	tests := []struct {
		email    string
		expected *models.UserDB
	}{
		{"test@example.com", &models.MockUsers()[0]},
		{"nonexistent@example.com", nil},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := db.FindUserByEmail(tt.email)

			if tt.expected == nil && result != nil {
				t.Errorf("expected nil, got %v", result)
				return
			}
			if tt.expected != nil && result == nil {
				t.Errorf("expected %v, got nil", tt.expected)
				return
			}
			if tt.expected != nil && result.ID != tt.expected.ID {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFindUserByID(t *testing.T) {
	db := NewMockDB()

	tests := []struct {
		id       string
		expected *models.UserDB
	}{
		{"user1", &models.MockUsers()[0]},
		{"nonexistent", nil},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			result := db.FindUserByID(tt.id)

			if tt.expected == nil && result != nil {
				t.Errorf("expected nil, got %v", result)
				return
			}
			if tt.expected != nil && result == nil {
				t.Errorf("expected %v, got nil", tt.expected)
				return
			}
			if tt.expected != nil && result.ID != tt.expected.ID {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFindMovies(t *testing.T) {
	db := NewMockDB()
	tests := []struct {
		offset   uint
		limit    uint
		expected []models.Movie
	}{
		{0, 2, db.movies[0:2]},
		{1, 2, db.movies[1:3]},
		{3, 2, db.movies[3:5]},
		{uint(len(db.movies)), 2, []models.Movie{}}, // Offset beyond range
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("offset=%d limit=%d", tt.offset, tt.limit), func(t *testing.T) {
			result := db.FindMovies(tt.offset, tt.limit)
			log.Println(result)

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d movies, got %d", len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i].ID != tt.expected[i].ID {
					t.Errorf("expected movie ID %s, got %s", tt.expected[i].ID, result[i].ID)
				}
			}
		})
	}
}
