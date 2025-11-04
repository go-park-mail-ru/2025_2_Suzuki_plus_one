package controller

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

//go:generate mockgen -source=contract.go -destination=./contract_mock.go -package=controller
type (
	GetMovieRecommendationsUsecase interface {
		Execute(context.Context, dto.GetMovieRecommendationsInput) (dto.GetMovieRecommendationsOutput, *dto.Error)
	}

	GetObjectUsecase interface {
		Execute(context.Context, dto.GetObjectInput) (dto.GetObjectOutput, *dto.Error)
	}

	PostAuthSignInUsecase interface {
		Execute(context.Context, dto.PostAuthSignInInput) (dto.PostAuthSignInOutput, *dto.Error)
	}

	GetAuthRefreshUsecase interface {
		Execute(context.Context, dto.GetAuthRefreshInput) (dto.GetAuthRefreshOutput, *dto.Error)
	}
)
