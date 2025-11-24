package auth

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/entity"
)

//go:generate mockgen -source=contract.go -destination=./mocks/contract_mock.go -package=mocks
type (
	TokenRepository interface {
		AddNewRefreshToken(ctx context.Context, userID uint, refreshToken string, expiresAt time.Time) error
		GetRefreshTokensForUser(ctx context.Context, userID uint) ([]entity.RefreshToken, error)
		RemoveRefreshToken(ctx context.Context, userID uint, refreshToken string) error
	}

	UserRepository interface {
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		CreateUser(ctx context.Context, email string, username string, passwordHash string) (uint, error)
	}
)
