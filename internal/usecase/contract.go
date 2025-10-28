package usecase

import (
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

type (
	MovieRepository interface {
		GetMovies(limit, offset uint) ([]entity.Movie, error)
	}
)