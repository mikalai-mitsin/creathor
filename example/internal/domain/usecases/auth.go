package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// AuthUseCase - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/auth.go . AuthUseCase
type AuthUseCase interface {
	CreateToken(ctx context.Context, login *models.Login) (*models.TokenPair, error)
	CreateTokenByUser(ctx context.Context, user *models.User) (*models.TokenPair, error)
	HasPermission(ctx context.Context, user *models.User, permission models.PermissionID) error
	HasObjectPermission(
		ctx context.Context,
		user *models.User,
		permission models.PermissionID,
		object any,
	) error
	RefreshToken(ctx context.Context, refresh models.Token) (*models.TokenPair, error)
	ValidateToken(ctx context.Context, access models.Token) error
	Auth(ctx context.Context, access models.Token) (*models.User, error)
}
