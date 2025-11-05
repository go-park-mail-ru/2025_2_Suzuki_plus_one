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
	movieRepo := NewMockMediaRepository(mockCtrl)
	times := 5

	// Media count times*2
	movieRepo.EXPECT().GetMediaCount(gomock.Any(), "movie").Return(times*2, nil).Times(1)
	// Rest is times because of the limit input
	movieRepo.EXPECT().GetMedia(gomock.Any(), gomock.Any()).Return(&entity.Media{}, nil).Times(times)
	movieRepo.EXPECT().GetMediaGenres(gomock.Any(), gomock.Any()).Return([]entity.Genre{}, nil).Times(times)
	movieRepo.EXPECT().GetMediaPostersKeys(gomock.Any(), gomock.Any()).Return([]string{"posters/hi.png"}, nil).Times(times)
	movieRepo.EXPECT().GetActorsByMediaID(gomock.Any(), gomock.Any()).Return([]entity.Actor{}, nil).Times(times)

	// Object repository mock
	objectRepo := NewMockObjectRepository(mockCtrl)
	// For each poster key, GetObject will be called
	objectRepo.EXPECT().GetPublicObject(gomock.Any(), gomock.Any(), "posters").Return(&entity.Object{
		URL: "http://example.com/poster.jpg",
	}, nil).Times(times)

	// Call usecase
	logger := logger.NewZapLogger(true)
	usecase := NewGetMovieRecommendationsUsecase(
		logger,
		movieRepo,
		NewGetMediaUseCase(
			logger,
			movieRepo,
			NewGetObjectUseCase(logger, objectRepo),
		),
	)
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
