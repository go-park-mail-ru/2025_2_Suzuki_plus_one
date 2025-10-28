package controller

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/dto"
)

type (
	GetMoviesUsecase interface {
		Execute(dto.GetMoviesInput) (dto.GetMoviesOutput, *dto.Error)
	}
)
