package http

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	. "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/usecase/http/mocks"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetAuthRefreshUsecase(t *testing.T) {
	// Init repository mock
	mockCtrl := gomock.NewController(t)
	authService := NewMockServiceAuthRepository(mockCtrl)
	common.InitJWT("secret", time.Hour, time.Hour*24*7)

	userID := uint(1)
	expiresAt := time.Hour
	token, genErr := common.GenerateToken(uint(userID), expiresAt)
	require.NoError(t, genErr)
	mockAccessToken, err := common.GenerateToken(userID, common.AccessTokenTTL)
	require.NoError(t, err)

	// Media count times*2
	authService.EXPECT().CallRefresh(gomock.Any(), token).Return(mockAccessToken, nil).Times(1)

	// Call usecase
	logger := logger.NewZapLogger(true)
	usecase := NewGetAuthRefreshUseCase(logger, authService)
	ctx := context.Background()
	output, dtoErr := usecase.Execute(ctx, dto.GetAuthRefreshInput{
		RefreshToken: token,
	})
	var emptyErr *dto.Error
	require.Equal(t, dtoErr, emptyErr)

	// Check output
	require.NotEmpty(t, output.AccessToken)
	userId, qerr := common.ValidateToken(output.AccessToken)
	require.Equal(t, qerr, nil)
	require.Equal(t, userID, userId, "User IDs must be equal")
}
