package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostAuthSignInUsecase struct {
	logger      logger.Logger
	authRepo    UserRepository
	sessionRepo SessionRepository
}

func NewPostAuthSignInUsecase(
	logger logger.Logger,
	authRepo UserRepository,
	sessionRepo SessionRepository,
) *PostAuthSignInUsecase {
	return &PostAuthSignInUsecase{
		logger:      logger,
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *PostAuthSignInUsecase) Execute(ctx context.Context, input dto.PostAuthSignInInput) (dto.PostAuthSignInOutput, *dto.Error) {
	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInParamsInvalid,
			err.Error(),
		)
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Check user credentials in repository and get user id
	user, err := uc.authRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInParamsInvalid,
			"invalid email or password",
		)
		uc.logger.Warn("Invalid email attempt", uc.logger.ToString("email", input.Email))
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Verify password with bcrypt
	if err := common.ValidateHashedPasswordBcrypt(user.PasswordHash, input.Password); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInParamsInvalid,
			"invalid email or password",
		)
		uc.logger.Warn("Invalid password attempt",
			uc.logger.ToString("password", input.Password),
			uc.logger.ToError(err))
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Add refresh refreshToken for user in repository
	refreshToken := generateRefreshToken(user.ID)
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // TODO: 7 days expiration, can be changed later if needed

	if err := uc.authRepo.AddNewRefreshToken(ctx, user.ID, refreshToken, expiresAt); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInNewRefreshTokenFailed,
			err.Error(),
		)
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Generate access token
	accessToken := generateAccessToken(user.ID)
	accessTokenExpiration := time.Hour // TODO: 1 hour expiration, can be changed later if needed

	// Create session in cache (Redis)
	if err := uc.sessionRepo.AddSession(ctx, user.ID, accessToken, accessTokenExpiration); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInAddSessionFailed,
			err.Error(),
		)
		return dto.PostAuthSignInOutput{}, &derr
	}

	return dto.PostAuthSignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateRefreshToken(userID uint) string {
	// Placeholder implementation for refresh token generation
	return fmt.Sprintf("refresh-token-for-user-%d", userID)
}

func generateAccessToken(userID uint) string {
	// Placeholder implementation for access token generation
	return fmt.Sprintf("access-token-for-user-%d", userID)
}
