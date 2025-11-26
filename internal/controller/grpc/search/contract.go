package search

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

type (
	SearchMediaUsecase interface {
		Execute(ctx context.Context, input dto.GetSearchInput) (dto.GetSearchOutput, *dto.Error)
	}

	SearchActorUsecase interface {
		Execute(ctx context.Context, input dto.GetSearchInput) (dto.GetSearchOutput, *dto.Error)
	}
)
