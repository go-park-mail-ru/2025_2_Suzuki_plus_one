package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	logpkg "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestGetMediaCount(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.ContextKeyRequestID, "1")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	// Prepare expected row
	rows := sqlmock.NewRows([]string{"count"}).AddRow(42)

	// Expect the query
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM media").WillReturnRows(rows)

	count, err := database.GetMediaCount(ctx, "movie")
	require.NoError(t, err)

	require.Equal(t, 42, count)

	// ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetMedia(t *testing.T) {
	ctx := context.WithValue(context.Background(), common.ContextKeyRequestID, "1")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	lg := logpkg.NewZapLogger(true)

	database := &DataBase{
		conn:    db,
		logger:  lg,
		context: ctx,
	}

	// Prepare expected row (match the SELECT in GetMovie)
	releaseDate, _ := time.Parse("2006-01-02", "2010-07-16")
	rows := sqlmock.NewRows([]string{
		"media_id",
		"media_type",
		"title",
		"description",
		"release_date",
		"rating",
		"duration_minutes",
		"age_rating",
		"country",
		"plot_summary",
	}).AddRow(
		1,
		"movie",
		"Inception",
		"A mind-bending thriller",
		releaseDate,
		8.8,
		148,
		13,
		"USA",
		"A thief who steals corporate secrets...",
	)

	// Expect the query. Use a partial regex to avoid whitespace/newline mismatch.
	mock.ExpectQuery(`FROM media\s+WHERE media_id = \$1`).WillReturnRows(rows)

	media, err := database.GetMediaByID(ctx, 1)
	require.NoError(t, err)

	require.Equal(t, uint(1), media.MediaID)
	require.Equal(t, "movie", media.MediaType)
	require.Equal(t, "Inception", media.Title)
	require.Equal(t, "A mind-bending thriller", media.Description)
	require.Equal(t, releaseDate, media.ReleaseDate)
	require.InDelta(t, 8.8, media.Rating, 0.001)
	require.Equal(t, 148, media.Duration)
	require.Equal(t, 13, media.AgeRating)
	require.Equal(t, "USA", media.Country)
	require.Equal(t, "A thief who steals corporate secrets...", media.PlotSummary)

	// ensure all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
