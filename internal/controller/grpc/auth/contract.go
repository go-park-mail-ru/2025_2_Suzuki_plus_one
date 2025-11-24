package auth

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

type (
	LoginUsecase interface {
		Execute(ctx context.Context, input dto.PostAuthSignInInput) (dto.PostAuthSignInOutput, *dto.Error)
	}

	RefreshUsecase interface {
		Execute(ctx context.Context, input dto.GetAuthRefreshInput) (dto.GetAuthRefreshOutput, *dto.Error)
	}

	LogoutUsecase interface {
		Execute(ctx context.Context, input dto.GetAuthSignOutInput) (dto.GetAuthSignOutOutput, *dto.Error)
	}

	CreateUserUsecase interface {
		Execute(ctx context.Context, input dto.PostAuthSignUpInput) (dto.PostAuthSignUpOutput, *dto.Error)
	}
)
