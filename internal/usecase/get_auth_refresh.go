package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetAuthRefreshUseCase struct {
	logger    logger.Logger
	tokenRepo TokenRepository
}

func NewGetAuthRefreshUseCase(logger logger.Logger, tokenRepo TokenRepository) *GetAuthRefreshUseCase {
	return &GetAuthRefreshUseCase{
		logger:    logger,
		tokenRepo: tokenRepo,
	}
}

func (u *GetAuthRefreshUseCase) Execute(
	ctx context.Context,
	input dto.GetAuthRefreshInput,
) (
	dto.GetAuthRefreshOutput,
	*dto.Error,
) {
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		return dto.GetAuthRefreshOutput{}, &derr
	}
	u.logger.Info("GetAuthRefreshUseCase called",
		u.logger.ToString("refresh_token", input.RefreshToken),
	)

	// Validate refresh token
	userID, err := common.ValidateToken(input.RefreshToken)

	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		return dto.GetAuthRefreshOutput{}, &derr
	}

	// Check for user existence
	tokens, err := u.tokenRepo.GetRefreshTokensForUser(ctx, userID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		return dto.GetAuthRefreshOutput{}, &derr
	}
	// Check if the provided refresh token exists for the user
	hasToken := false
	for _, t := range tokens {
		if t.Token == input.RefreshToken {
			hasToken = true
			break
		}
	}
	if !hasToken {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			"refresh token not found for user",
		)
		return dto.GetAuthRefreshOutput{}, &derr
	}

	// Generate new access token
	accessToken, err := common.GenerateToken(userID, common.AccessTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		return dto.GetAuthRefreshOutput{}, &derr
	}

	return dto.GetAuthRefreshOutput{AccessToken: accessToken}, nil
}
