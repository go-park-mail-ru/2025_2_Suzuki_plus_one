package auth

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type CreateUserUseCase struct {
	logger    logger.Logger
	userRepo  UserRepository
	tokenRepo TokenRepository
}

func NewCreateUserUsecase(l logger.Logger, ur UserRepository, tr TokenRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		logger:    l,
		userRepo:  ur,
		tokenRepo: tr,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.PostAuthSignUpInput) (dto.PostAuthSignUpOutput, *dto.Error) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"auth/usecase/post_auth_signup",
			entity.ErrPostAuthSignUpParamsInvalid,
			err.Error(),
		)
		log.Error("Invalid sign up input parameters", log.ToError(err))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Check for given email if user already exists
	_, err := uc.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil {
		derr := dto.NewError(
			"auth/usecase/post_auth_signup",
			entity.ErrPostAuthSignUpAlreadyExists,
			"user with given email already exists",
		)
		log.Warn("User already exists", log.ToString("email", input.Email))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Try to create new user in repository
	hashedPassword, err := common.HashPasswordBcrypt(input.Password)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/post_auth_signup",
			entity.ErrPostAuthSignUpParamsInvalid,
			"failed to hash password",
		)
		log.Error("Failed to hash password", log.ToError(err))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	userID, err := uc.userRepo.CreateUser(ctx, input.Email, input.Username, hashedPassword)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/post_auth_signup",
			entity.ErrPostAuthSignUpParamsInvalid,
			"failed to create new user: "+err.Error(),
		)
		log.Error("Failed to create new user", log.ToError(err))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Generate refresh token
	refreshToken, err := common.GenerateToken(userID, common.RefreshTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/post_auth_signup",
			entity.ErrPostAuthSignInNewRefreshTokenFailed,
			err.Error(),
		)
		log.Error("Failed to generate refresh token", log.ToError(err))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Add refresh refreshToken for user in repository
	expiration := time.Now().Add(common.RefreshTokenTTL)
	if err := uc.tokenRepo.AddNewRefreshToken(ctx, userID, refreshToken, expiration); err != nil {
		derr := dto.NewError(
			"auth/usecase/post_auth_signup",
			entity.ErrPostAuthSignInNewRefreshTokenFailed,
			err.Error(),
		)
		log.Error("Failed to add new refresh token", log.ToError(err))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Generate access token
	accessToken, err := common.GenerateToken(userID, common.AccessTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"auth/usecase/post_auth_signup",
			entity.ErrPostAuthSignInNewAccessTokenFailed,
			err.Error(),
		)
		log.Error("Failed to generate access token", log.ToError(err))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	return dto.PostAuthSignUpOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
