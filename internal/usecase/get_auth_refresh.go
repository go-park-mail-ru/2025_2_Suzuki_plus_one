package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetAuthRefreshUseCase struct {
	logger      logger.Logger
	tokenRepo   TokenRepository
	sessionRepo SessionRepository
}

func NewGetAuthRefreshUseCase(logger logger.Logger, tokenRepo TokenRepository, sessionRepo SessionRepository) *GetAuthRefreshUseCase {
	return &GetAuthRefreshUseCase{
		logger:      logger,
		tokenRepo:   tokenRepo,
		sessionRepo: sessionRepo,
	}
}

func (u *GetAuthRefreshUseCase) Execute(
	ctx context.Context,
	input dto.GetAuthRefreshInput,
) (
	dto.GetAuthRefreshOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(u.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			err.Error(),
		)
		log.Error("GetAuthRefreshUseCase failed on validation", log.ToError(err))
		return dto.GetAuthRefreshOutput{}, &derr
	}
	log.Info("GetAuthRefreshUseCase called",
		log.ToString("refresh_token", input.RefreshToken),
	)

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

	// TODO: remove old access tokens in redis

	// Check for user existence
	tokens, err := u.tokenRepo.GetRefreshTokensForUser(ctx, userID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
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
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			"refresh token not found for user",
		)
		log.Error("Refresh token not found for user", log.ToInt("user_id", int(userID)))
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
		log.Error("Failed to generate access token", log.ToError(err))
		return dto.GetAuthRefreshOutput{}, &derr
	}

	// Add new access token to session repository
	err = u.sessionRepo.AddSession(ctx, userID, accessToken, common.AccessTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshInvalidParams,
			"failed to add session: "+err.Error(),
		)
		log.Error("Failed to add session", log.ToError(err))
		return dto.GetAuthRefreshOutput{}, &derr
	}

	return dto.GetAuthRefreshOutput{AccessToken: accessToken}, nil
}
