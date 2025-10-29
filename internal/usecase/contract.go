package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

//go:generate mockgen -source=contract.go -destination=./contract_mock.go -package=usecase
type (
	MovieRepository interface {
		GetMovies(ctx context.Context, limit, offset uint) ([]entity.Movie, error)
	}
)
