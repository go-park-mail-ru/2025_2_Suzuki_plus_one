package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type PostAuthSignUpUsecase struct {
	logger      logger.Logger
	authService ServiceAuthRepository
}

func NewPostAuthSignUpUsecase(
	logger logger.Logger,
	authService ServiceAuthRepository,
) *PostAuthSignUpUsecase {
	return &PostAuthSignUpUsecase{
		logger:      logger,
		authService: authService,
	}
}

func (uc *PostAuthSignUpUsecase) Execute(
	ctx context.Context,
	input dto.PostAuthSignUpInput,
) (
	dto.PostAuthSignUpOutput,
	*dto.Error,
) {
	// Bind logger with request ID
	log := logger.LoggerWithKey(uc.logger, ctx, common.ContextKeyRequestID)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signup",
			entity.ErrPostAuthSignUpParamsInvalid,
			err.Error(),
		)
		log.Error("Invalid sign up input parameters", log.ToError(err))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Call Auth service to create new user
	_, accessToken, refreshToken, err := uc.authService.CallCreateUser(ctx, input.Email, input.Username, input.Password)
	if err != nil {
		derr := dto.NewError(
			"usecase/post_auth_signup",
			entity.ErrPostAuthSignUpAuthServiceFailed,
			"authentication service failed: "+err.Error(),
		)
		log.Error("Authentication service call failed", log.ToError(err))
		return dto.PostAuthSignUpOutput{}, &derr
	}

	// Return output DTO
	return dto.PostAuthSignUpOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
