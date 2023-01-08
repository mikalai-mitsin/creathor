package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth_mock.go github.com/018bf/example/internal/domain/interceptors AuthInterceptor

type AuthInterceptor interface {
	CreateToken(
		ctx context.Context,
		login *models.Login,
	) (*models.TokenPair, error)
	RefreshToken(
		ctx context.Context,
		refresh models.Token,
	) (*models.TokenPair, error)
	Auth(
		ctx context.Context,
		access models.Token,
	) (*models.User, error)
	ValidateToken(
		ctx context.Context,
		access models.Token,
	) error
}
