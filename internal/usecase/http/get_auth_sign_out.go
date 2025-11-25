package http

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/common"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type GetAuthSignOutUsecase struct {
	logger      logger.Logger
	authService ServiceAuthRepository
}

func NewGetAuthSignOutUsecase(
	logger logger.Logger,
	authService ServiceAuthRepository,
) *GetAuthSignOutUsecase {
	return &GetAuthSignOutUsecase{
		logger:      logger,
		authService: authService,
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
	log.Debug("GetAuthSignOutUsecase called",
		log.ToString("refresh_token", input.RefreshToken),
	)

	// Validate input
	if err := dto.ValidateStruct(input); err != nil {
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutInvalidParams,
			err.Error(),
		)
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Call Auth service to sign out
	if err := uc.authService.CallLogout(ctx, input.RefreshToken); err != nil {
		derr := dto.NewError(
			"usecase/get_auth_signout",
			entity.ErrGetAuthSignOutAuthServiceFailed,
			"authentication service failed: "+err.Error(),
		)
		log.Error("GetAuthSignOutUsecase failed on auth service call", log.ToError(err))
		return dto.GetAuthSignOutOutput{}, &derr
	}

	// Return output
	return dto.GetAuthSignOutOutput{}, nil
}
