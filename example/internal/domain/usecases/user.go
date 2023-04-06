package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// UserUseCase - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/user.go . UserUseCase
type UserUseCase interface {
	Get(ctx context.Context, id models.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, filter *models.UserFilter) ([]*models.User, uint64, error)
	Create(ctx context.Context, create *models.UserCreate) (*models.User, error)
	Update(ctx context.Context, update *models.UserUpdate) (*models.User, error)
	Delete(ctx context.Context, id models.UUID) error
}
