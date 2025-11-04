package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetMovieRecommendationsUsecase(t *testing.T) {
	// Init repository mock
	mockCtrl := gomock.NewController(t)
	movieRepo := NewMockMovieRepository(mockCtrl)
	times := 5

	// Media count times*2
	movieRepo.EXPECT().GetMediaCount(gomock.Any(), "movie").Return(times*2, nil).Times(1)
	// Rest is times because of the limit input
	movieRepo.EXPECT().GetMedia(gomock.Any(), gomock.Any()).Return(&entity.Media{}, nil).Times(times)
	movieRepo.EXPECT().GetMediaGenres(gomock.Any(), gomock.Any()).Return([]entity.Genre{}, nil).Times(times)
	movieRepo.EXPECT().GetMediaPostersLinks(gomock.Any(), gomock.Any()).Return([]string{}, nil).Times(times)

	// Call usecase
	logger := logger.NewZapLogger(true)
	usecase := NewGetMovieRecommendationsUsecase(logger, movieRepo)
	ctx := context.Background()
	output, err := usecase.Execute(ctx, dto.GetMovieRecommendationsInput{
		Limit:  uint(times),
		Offset: 0,
	})
	var emptyErr *dto.Error
	require.Equal(t, err, emptyErr)

	// Compare by length
	require.Equal(t, len(output.Movies), times)
}
