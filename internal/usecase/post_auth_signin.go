package usecase

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostAuthSignInUsecase struct {
	logger      logger.Logger
	userRepo    UserRepository
	authRepo    TokenRepository
	sessionRepo SessionRepository
}

func NewPostAuthSignInUsecase(
	logger logger.Logger,
	userRepo UserRepository,
	authRepo TokenRepository,
	sessionRepo SessionRepository,
) *PostAuthSignInUsecase {
	return &PostAuthSignInUsecase{
		logger:      logger,
		userRepo:    userRepo,
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *PostAuthSignInUsecase) Execute(
	ctx context.Context,
	input dto.PostAuthSignInInput,
) (
	dto.PostAuthSignInOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

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
	user, err := uc.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInParamsInvalid,
			"invalid email or password: "+err.Error(),
		)
		log.Warn("Invalid email attempt", log.ToString("email", input.Email))
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Verify password with bcrypt
	if err := common.ValidateHashedPasswordBcrypt(user.PasswordHash, input.Password); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInParamsInvalid,
			"invalid email or password",
		)
		log.Warn(
			"Invalid password attempt",
			log.ToString("password", input.Password),
			log.ToError(err),
		)
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Generate refresh token
	refreshToken, err := common.GenerateToken(user.ID, common.RefreshTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInNewRefreshTokenFailed,
			err.Error(),
		)
		log.Error("Failed to generate refresh token", log.ToError(err))
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Add refresh refreshToken for user in repository
	expiration := time.Now().Add(common.RefreshTokenTTL)
	if err := uc.authRepo.AddNewRefreshToken(ctx, user.ID, refreshToken, expiration); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInNewRefreshTokenFailed,
			err.Error(),
		)
		log.Error("Failed to add new refresh token", log.ToError(err))
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Generate access token
	accessToken, err := common.GenerateToken(user.ID, common.AccessTokenTTL)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInNewAccessTokenFailed,
			err.Error(),
		)
		log.Error("Failed to generate access token", log.ToError(err))
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Create session in cache (Redis)
	if err := uc.sessionRepo.AddSession(ctx, user.ID, accessToken, common.AccessTokenTTL); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInAddSessionFailed,
			err.Error(),
		)
		log.Error("Failed to add session", log.ToError(err))
		return dto.PostAuthSignInOutput{}, &derr
	}

	return dto.PostAuthSignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
