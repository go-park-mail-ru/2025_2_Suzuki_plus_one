package auth

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type LogoutUseCase struct {
	logger    logger.Logger
	userRepo  UserRepository
	tokenRepo TokenRepository
}

func NewLogoutUsecase(l logger.Logger, ur UserRepository, tr TokenRepository) *LogoutUseCase {
	return &LogoutUseCase{
		logger:    l,
		userRepo:  ur,
		tokenRepo: tr,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, input dto.GetAuthSignOutInput) (dto.GetAuthSignOutOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			err.Error(),
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Get refresh token
	userIDRefresh, err := common.ValidateToken(input.RefreshToken)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"invalid refresh token: "+err.Error(),
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Get access token user ID
	userIDAccess, err := common.ValidateToken(input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"invalid access token: "+err.Error(),
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Compare user IDs from access and refresh tokens
	if userIDAccess != userIDRefresh {
		log.Error("Access token user ID does not match refresh token user ID")
		derr := dto.NewError(
			"auth/usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"access token user ID does not match refresh token user ID",
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Remove refresh token from database
	if err := uc.tokenRepo.RemoveRefreshToken(ctx, userIDAccess, input.RefreshToken); err != nil {
		derr := dto.NewError(
			"auth/usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			"failed to remove refresh token: "+err.Error(),
		)
		log.Error("Failed to remove refresh token", log.ToError(err))
		return dto.GetAuthSignOutOutput{}, &derr
	}

	return dto.GetAuthSignOutOutput{}, nil
}
