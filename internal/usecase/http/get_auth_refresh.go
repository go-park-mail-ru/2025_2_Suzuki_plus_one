package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetAuthRefreshUseCase struct {
	logger      logger.Logger
	authService ServiceAuthRepository
}

func NewGetAuthRefreshUseCase(logger logger.Logger, authService ServiceAuthRepository) *GetAuthRefreshUseCase {
	return &GetAuthRefreshUseCase{
		logger:      logger,
		authService: authService,
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
	log.Debug("GetAuthRefreshUseCase called",
		log.ToString("refresh_token", input.RefreshToken),
	)

	// Call Auth service to refresh tokens
	accessToken, err := u.authService.CallRefresh(ctx, input.RefreshToken)
	if err != nil {
		derr := dto.NewError(
			"usecase/get_auth_refresh",
			entity.ErrGetAuthRefreshAuthServiceFailed,
			"authentication service failed: "+err.Error(),
		)
		log.Error("GetAuthRefreshUseCase failed on auth service call", log.ToError(err))
		return dto.GetAuthRefreshOutput{}, &derr
	}

	return dto.GetAuthRefreshOutput{AccessToken: accessToken}, nil
}
