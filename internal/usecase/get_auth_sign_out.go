package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetAuthSignOutUsecase struct {
	logger      logger.Logger
	authRepo    TokenRepository
	sessionRepo SessionRepository
}

func NewGetAuthSignOutUsecase(
	logger logger.Logger,
	authRepo TokenRepository,
	sessionRepo SessionRepository,
) *GetAuthSignOutUsecase {
	return &GetAuthSignOutUsecase{
		logger:      logger,
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *GetAuthSignOutUsecase) Execute(
	ctx context.Context,
	input dto.GetAuthSignOutInput,
) (
	dto.GetAuthSignOutOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			err.Error(),
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}
	// Get refresh token
	userIDRefresh, err := common.ValidateToken(input.RefreshToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"invalid refresh token: "+err.Error(),
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Get access token user ID
	userIDAccess, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"invalid access token: "+err.Error(),
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Compare user IDs from access and refresh tokens
	if userIDAccess != userIDRefresh {
		log.Error("Access token user ID does not match refresh token user ID")
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"access token user ID does not match refresh token user ID",
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Check if user authorized
	userIDCached, err := uc.sessionRepo.GetUserIDByAccessToken(ctx, input.AccessToken)
	if err != nil || userIDCached != userIDAccess {
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"user is not authorized: "+err.Error(),
		)
		log.Error("User is not authorized", log.ToError(err))
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Remove refresh token from database
	if err := uc.authRepo.RemoveRefreshToken(ctx, userIDAccess, input.RefreshToken); err != nil {
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"failed to remove refresh token: "+err.Error(),
		)
		log.Error("Failed to remove refresh token", log.ToError(err))
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Remove access token from Redis
	if err := uc.sessionRepo.DeleteSession(ctx, userIDAccess, input.AccessToken); err != nil {
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"failed to delete session: "+err.Error(),
		)
		log.Error("Failed to delete session", log.ToError(err))
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Return output
	return dto.GetAuthSignOutOutput{}, nil
}
