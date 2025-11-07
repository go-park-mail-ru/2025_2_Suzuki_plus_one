package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase/mocks"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetObjectUsecase_Execute(t *testing.T) {
	// Init repository mock
	mockCtrl := gomock.NewController(t)
	objectRepo := NewMockObjectRepository(mockCtrl)

	// Define input and expected output
	input := dto.GetObjectInput{
		Key:        "1",
		BucketName: "posters",
	}
	expectedObject := &entity.URL{
		URL: "http://example.com/object1.jpg",
	}

	// Note: "posters" is a public bucket, so GeneratePublicURL will be called
	objectRepo.EXPECT().
		GeneratePublicURL(gomock.Any(), input.BucketName, input.Key).
		Return(expectedObject, nil).
		Times(1)

	// For now only "medias" bucket is private and uses GeneratePresignedURL with linkAliveDuration
	// Set up expectations
	// objectRepo.EXPECT().
	// 	GeneratePresignedURL(gomock.Any(), input.Key, input.BucketName, gomock.Any()).
	// 	Return(expectedObject, nil).
	// 	Times(1)

	// Call usecase
	logger := logger.NewZapLogger(true)
	usecase := NewGetObjectUseCase(logger, objectRepo)
	ctx := context.Background()
	output, err := usecase.Execute(ctx, input)
	var emptyErr *dto.Error
	require.Equal(t, err, emptyErr)

	// Compare output
	require.Equal(t, output.URL, expectedObject.URL)
}
