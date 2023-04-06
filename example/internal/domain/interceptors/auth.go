package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// AuthInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth.go . AuthInterceptor
type AuthInterceptor interface {
	CreateToken(ctx context.Context, login *models.Login) (*models.TokenPair, error)
	RefreshToken(ctx context.Context, refresh models.Token) (*models.TokenPair, error)
	Auth(ctx context.Context, access models.Token) (*models.User, error)
	ValidateToken(ctx context.Context, access models.Token) error
}
