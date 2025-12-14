package auth

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type RefreshUseCase struct {
	logger    logger.Logger
	userRepo  UserRepository
	tokenRepo TokenRepository
}

func NewRefreshUsecase(l logger.Logger, ur UserRepository, tr TokenRepository) *RefreshUseCase {
	return &RefreshUseCase{
		logger:    l,
		userRepo:  ur,
		tokenRepo: tr,
	}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, input dto.GetAuthRefreshInput) (dto.GetAuthRefreshOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)
	log.Debug("auth/usecase/get_auth_refresh: called")

	// 1. Validation
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		return dto.GetAuthRefreshOutput{}, &derr
	}

	// Validate refresh token
	userID, err := common.ValidateToken(input.RefreshToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		log.Error("GetAuthRefreshUseCase failed on refresh token validation", log.ToError(err))
		return dto.GetAuthRefreshOutput{}, &derr
	}

	// 2. Compare refresh token with stored tokens
	// Get active refresh tokens for the user
	tokens, err := uc.tokenRepo.GetRefreshTokensForUser(ctx, userID)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		log.Error("GetAuthRefreshUseCase failed on getting refresh tokens for user",
			log.ToInt("user_id", int(userID)),
			log.ToError(err),
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
			"auth/usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			"refresh token not found for user",
		)
		log.Error("Refresh token not found for user", log.ToInt("user_id", int(userID)))
		return dto.GetAuthRefreshOutput{}, &derr
	}

	// 3. Generate new access token
	accessToken, err := common.GenerateToken(userID, common.AccessTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		log.Error("Failed to generate access token", log.ToError(err))
		return dto.GetAuthRefreshOutput{}, &derr
	}

	return dto.GetAuthRefreshOutput{
		AccessToken: accessToken,
	}, nil
}
