package controller

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

//go:generate mockgen -source=contract.go -destination=./mocks/contract_mock.go -package=mocks
type (
	GetMovieRecommendationsUseCase interface {
		Execute(context.Context, dto.GetMovieRecommendationsInput) (dto.GetMovieRecommendationsOutput, *dto.Error)
	}

	GetObjectUseCase interface {
		Execute(context.Context, dto.GetObjectInput) (dto.GetObjectOutput, *dto.Error)
	}

	PostAuthSignInUseCase interface {
		Execute(context.Context, dto.PostAuthSignInInput) (dto.PostAuthSignInOutput, *dto.Error)
	}

	GetAuthRefreshUseCase interface {
		Execute(context.Context, dto.GetAuthRefreshInput) (dto.GetAuthRefreshOutput, *dto.Error)
	}

	PostAuthSignUpUseCase interface {
		Execute(context.Context, dto.PostAuthSignUpInput) (dto.PostAuthSignUpOutput, *dto.Error)
	}

	GetAuthSignOutUseCase interface {
		Execute(context.Context, dto.GetAuthSignOutInput) (dto.GetAuthSignOutOutput, *dto.Error)
	}

	GetUserMeUseCase interface {
		Execute(context.Context, dto.GetUserMeInput) (dto.GetUserMeOutput, *dto.Error)
	}

	GetActorUseCase interface {
		Execute(context.Context, dto.GetActorInput) (dto.GetActorOutput, *dto.Error)
	}

	GetMediaUseCase interface {
		Execute(context.Context, dto.GetMediaInput) (dto.GetMediaOutput, *dto.Error)
	}
)
