package usecase

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostAuthSignUpUsecase struct {
	logger      logger.Logger
	userRepo    UserRepository
	authRepo    TokenRepository
	sessionRepo SessionRepository
}

func NewPostAuthSignUpUsecase(
	logger logger.Logger,
	userRepo UserRepository,
	authRepo TokenRepository,
	sessionRepo SessionRepository,
) *PostAuthSignUpUsecase {
	return &PostAuthSignUpUsecase{
		logger:      logger,
		userRepo:    userRepo,
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *PostAuthSignUpUsecase) Execute(
	ctx context.Context,
	input dto.PostAuthSignUpInput,
) (
	dto.PostAuthSignUpOutput,
	*dto.Error,
) {
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signup",
			entity.ErrPostAuthSignUpParamsInvalid,
			err.Error(),
		)
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Check for given email if user already exists
	_, err := uc.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil {
		derr := dto.NewError(
			"usecase/post_auth_signup",
			entity.ErrPostAuthSignUpAlreadyExists,
			"user with given email already exists",
		)
		uc.logger.Warn("User already exists", uc.logger.ToString("email", input.Email))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Try to create new user in repository
	hashedPassword, err := common.HashPasswordBcrypt(input.Password)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signup",
			entity.ErrPostAuthSignUpParamsInvalid,
			"failed to hash password",
		)
		return dto.PostAuthSignUpOutput{}, &derr
	}

	newUser := entity.User{
		Email:        input.Email,
		PasswordHash: hashedPassword,
		Username:     input.Username,
	}
	userID, err := uc.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signup",
			entity.ErrPostAuthSignUpParamsInvalid,
			"failed to create new user: "+err.Error(),
		)
		return dto.PostAuthSignUpOutput{}, &derr
	}
	newUser.ID = userID

	// Generate refresh token
	refreshToken, err := common.GenerateToken(newUser.ID, common.RefreshTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInNewRefreshTokenFailed,
			err.Error(),
		)
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Add refresh refreshToken for user in repository
	expiration := time.Now().Add(common.RefreshTokenTTL)
	if err := uc.authRepo.AddNewRefreshToken(ctx, newUser.ID, refreshToken, expiration); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInNewRefreshTokenFailed,
			err.Error(),
		)
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Generate access token
	accessToken, err := common.GenerateToken(newUser.ID, common.AccessTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInNewAccessTokenFailed,
			err.Error(),
		)
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Create session in cache (Redis)
	if err := uc.sessionRepo.AddSession(ctx, newUser.ID, accessToken, common.AccessTokenTTL); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInAddSessionFailed,
			err.Error(),
		)
		return dto.PostAuthSignUpOutput{}, &derr
	}

	return dto.PostAuthSignUpOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
