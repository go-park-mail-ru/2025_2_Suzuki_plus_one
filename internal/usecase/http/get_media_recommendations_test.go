package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase/http/mocks"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetMediaRecommendationsUsecase(t *testing.T) {
	// Init repository mock
	mockCtrl := gomock.NewController(t)
	mediaRepo := NewMockMediaRepository(mockCtrl)
	times := 5

	// Media count times*2
	mediaRepo.EXPECT().GetMediaCount(gomock.Any(), "movie").Return(times*2, nil).Times(1)
	// Rest is times because of the limit input
	mediaRepo.EXPECT().GetMediaByID(gomock.Any(), gomock.Any()).Return(&entity.Media{}, nil).Times(times)
	mediaRepo.EXPECT().GetMediaGenres(gomock.Any(), gomock.Any()).Return([]entity.Genre{}, nil).Times(times)
	// Adjusted to return []entity.S3Key (matches mock signature)
	mediaRepo.EXPECT().GetMediaPostersKeys(gomock.Any(), gomock.Any()).Return([]entity.S3Key{
		{Key: "posters/hi.png", BucketName: "posters"},
	}, nil).Times(times)
	mediaRepo.EXPECT().GetMediaTrailersKeys(gomock.Any(), gomock.Any()).Return([]entity.S3Key{
		{Key: "trailers/hi.mp4", BucketName: "trailers"},
	}, nil).Times(times)
	mediaRepo.EXPECT().GetMediaSortedByName(gomock.Any(), uint(times), uint(0), "movie").Return([]uint{1, 2, 3, 4, 5}, nil).Times(1)

	// Object repository mock
	objectRepo := NewMockObjectRepository(mockCtrl)
	// For each poster key, GetObjectURL will be called
	objectRepo.EXPECT().GeneratePublicURL(gomock.Any(), gomock.Any(), "posters/hi.png").Return(&entity.URL{
		URL: "http://example.com/poster.jpg",
	}, nil).Times(times)
	objectRepo.EXPECT().GeneratePublicURL(gomock.Any(), gomock.Any(), "trailers/hi.mp4").Return(&entity.URL{
		URL: "http://example.com/trailer.mp4",
	}, nil).Times(times)

	// Call usecase
	logger := logger.NewZapLogger(true)
	getMediaUseCase := NewGetMediaUseCase(
		logger,
		mediaRepo,
		NewGetObjectUseCase(logger, objectRepo),
	)
	usecase := NewGetMediaRecommendationsUsecase(
		logger,
		mediaRepo,
		getMediaUseCase,
	)
	ctx := context.Background()
	output, err := usecase.Execute(ctx, dto.GetMediaRecommendationsInput{
		Limit:  uint(times),
		Offset: 0,
		Type:   "movie",
	})
	var emptyErr *dto.Error
	require.Equal(t, err, emptyErr)

	// Compare by length
	require.Equal(t, len(output.Movies), times)
}
