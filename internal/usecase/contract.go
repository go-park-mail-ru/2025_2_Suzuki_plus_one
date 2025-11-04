package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

//go:generate mockgen -source=contract.go -destination=./contract_mock.go -package=usecase
type (
	Repository interface {
		Connect() error
		Close() error
	}
	MovieRepository interface {
		GetMediaCount(ctx context.Context, media_type string) (int, error)
		GetMedia(ctx context.Context, media_id uint) (*entity.Media, error)
		GetMediaGenres(ctx context.Context, media_id uint) ([]entity.Genre, error)
		GetMediaPostersLinks(ctx context.Context, media_id uint) ([]string, error)
	}
)
