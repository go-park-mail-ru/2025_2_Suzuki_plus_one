package entity

import (
	"fmt"
	"time"
)

type Media struct {
	MediaID     uint      `json:"media_id"`
	MediaType   string    `json:"media_type" validate:"required"`
	Title       string    `json:"title" validate:"required,max=256"`
	Description string    `json:"description,omitempty" validate:"max=2048"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	Rating      float64   `json:"rating,omitempty"`
	Duration    int       `json:"duration_minutes,omitempty" validate:"min=0"`
	AgeRating   int       `json:"age_rating,omitempty" validate:"min=0"`
	Country     string    `json:"country,omitempty" validate:"max=64"`
	PlotSummary string    `json:"plot_summary,omitempty" validate:"max=4096"`
}

func NewMedia(
	mediaID uint,
	mediaType string,
	title string,
	description string,
	releaseDate time.Time,
	rating float64,
	duration int,
	ageRating int,
	country string,
	plotSummary string,
) (Media, error) {
	m := Media{
		MediaID:     mediaID,
		MediaType:   mediaType,
		Title:       title,
		Description: description,
		ReleaseDate: releaseDate,
		Rating:      rating,
		Duration:    duration,
		AgeRating:   ageRating,
		Country:     country,
		PlotSummary: plotSummary,
	}
	// Validate the media entity
	if err := validate.Struct(m); err != nil {
		return Media{}, fmt.Errorf("invalid media entity validation: %w", err)
	}
	return m, nil
}

type Episode struct {
	EpisodeID     uint `json:"episode_id"` // References Media.MediaID
	SeriesID      uint `json:"series_id"`  // References Media.MediaID
	SeasonNumber  int  `json:"season_number" validate:"min=1"`
	EpisodeNumber int  `json:"episode_number" validate:"min=1"`
}
