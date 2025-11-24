package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostAuthSignInUsecase struct {
	logger      logger.Logger
	authService ServiceAuthRepository
}

func NewPostAuthSignInUsecase(
	logger logger.Logger,
	authService ServiceAuthRepository,
) *PostAuthSignInUsecase {
	return &PostAuthSignInUsecase{
		logger:      logger,
		authService: authService,
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
	log.Debug("Executing LoginUseCase with input: ", log.ToString("email", input.Email))

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInParamsInvalid,
			err.Error(),
		)
		log.Error("Invalid sign in input parameters", log.ToError(err))
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Call Auth microservice to verify user credentials and get tokens
	accessToken, refreshToken, err := uc.authService.CallLogin(ctx, input.Email, input.Password)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signin",
			entity.ErrPostAuthSignInAuthServiceFailed,
			"authentication service failed: "+err.Error(),
		)
		log.Error("Authentication service call failed", log.ToError(err))
		return dto.PostAuthSignInOutput{}, &derr
	}

	// Return output DTO
	return dto.PostAuthSignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
