package search

import (
	"context"
)

//go:generate mockgen -source=contract.go -destination=./mocks/contract_mock.go -package=mocks
type (
	MediaRepository interface {
		SearchMedia(ctx context.Context, query string, limit, offset uint) ([]uint, error)
	}

	ActorRepository interface {
		SearchActor(ctx context.Context, query string, limit, offset uint) ([]uint, error)
	}
)
