package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestGetMovieRecommendations(t *testing.T) {
	// Set up context with request ID
	ctx := context.WithValue(context.Background(), common.RequestIDContextKey, "1")

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "failed to open sqlmock database")
	defer db.Close()

	logger := logger.NewZapLogger(true)
	database := &DataBase{
		conn:    db,
		logger:  logger,
		context: ctx,
	}
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM movies").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	// Set up your mock expectations
	mock.ExpectQuery("SELECT id, title, year FROM movies ORDER BY RANDOM\\(\\) OFFSET \\$1 LIMIT \\$2").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "year"}).
			AddRow(1, "Inception", 2010).AddRow(2, "The Dark Knight", 2008))

	// Test the GetMovieRecommendations method
	recommendations, err := database.GetMovieRecommendations(ctx, 1, 10)
	require.NoError(t, err, "error was not expected while getting movie recommendations")
	logger.Info("Recommendations:", "data", recommendations)

	require.Equal(t, 3, len(recommendations), "expected 2 movie recommendations")
}
