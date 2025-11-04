package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetUserMeUseCase struct {
	logger      logger.Logger
	userRepo    UserRepository
	sessionRepo SessionRepository
}

func NewGetUserMeUseCase(
	logger logger.Logger,
	userRepo UserRepository,
	sessionRepo SessionRepository,
) *GetUserMeUseCase {
	return &GetUserMeUseCase{
		logger:      logger,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *GetUserMeUseCase) Execute(ctx context.Context, input dto.GetUserMeInput) (dto.GetUserMeOutput, *dto.Error) {
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrGetUserMeParamsInvalid,
			err.Error(),
		)
		return dto.GetUserMeOutput{}, &derr
	}

	// Get session by access token
	userID, err := uc.sessionRepo.GetUserIDByToken(ctx, input.AccessToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrGetUserMeSessionNotFound,
			err.Error(),
		)
		uc.logger.Error("Failed to get user ID by token", uc.logger.ToError(err))
		return dto.GetUserMeOutput{}, &derr
	}

	// Compare user ID from session with requested user ID
	userIDToken, err := common.ValidateToken(input.AccessToken)
	if err != nil || userID != userIDToken {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrGetUserMeSessionNotFound,
			"access token is invalid",
		)
		uc.logger.Error("Failed to validate access token", uc.logger.ToError(err))
		return dto.GetUserMeOutput{}, &derr
	}

	// Get user info
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_user_me",
			entity.ErrUserNotFound,
			err.Error(),
		)
		uc.logger.Error("Failed to get user by ID", uc.logger.ToError(err))
		return dto.GetUserMeOutput{}, &derr
	}

	return dto.GetUserMeOutput{User: *user}, nil
}
