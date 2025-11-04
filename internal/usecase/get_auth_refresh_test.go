package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestGetAuthRefreshUsecase(t *testing.T) {
	// Init repository mock
	mockCtrl := gomock.NewController(t)
	movieRepo := NewMockTokenRepository(mockCtrl)
	common.InitJWT("secret", time.Hour, time.Hour*24*7)

	userID := uint(1)
	refreshTokens := make([]entity.RefreshToken, 1)
	expiresAt := time.Hour
	token, genErr := common.GenerateToken(uint(userID), expiresAt)
	require.NoError(t, genErr)
	refreshTokens = append(refreshTokens,
		entity.RefreshToken{
			ID:        1,
			ExpiresAt: time.Now().Add(expiresAt),
			UserID:    userID,
			Token:     token,
		},
	)
	// Media count times*2
	movieRepo.EXPECT().GetRefreshTokensForUser(gomock.Any(), userID).Return(refreshTokens, nil).Times(1)

	// Call usecase
	logger := logger.NewZapLogger(true)
	usecase := NewGetAuthRefreshUseCase(logger, movieRepo)
	ctx := context.Background()
	output, err := usecase.Execute(ctx, dto.GetAuthRefreshInput{
		RefreshToken: token,
	})
	var emptyErr *dto.Error
	require.Equal(t, err, emptyErr)

	// Check output
	require.NotEmpty(t, output.AccessToken)
	userId, qerr := common.ValidateToken(output.AccessToken)
	require.Equal(t, qerr, nil)
	require.Equal(t, userID, userId, "User IDs must be equal")
}
