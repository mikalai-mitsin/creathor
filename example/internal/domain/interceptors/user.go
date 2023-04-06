package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// UserInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/user.go . UserInterceptor
type UserInterceptor interface {
	Get(ctx context.Context, id models.UUID, requestUser *models.User) (*models.User, error)
	List(
		ctx context.Context,
		filter *models.UserFilter,
		requestUser *models.User,
	) ([]*models.User, uint64, error)
	Create(
		ctx context.Context,
		create *models.UserCreate,
		requestUser *models.User,
	) (*models.User, error)
	Update(
		ctx context.Context,
		update *models.UserUpdate,
		requestUser *models.User,
	) (*models.User, error)
	Delete(ctx context.Context, id models.UUID, requestUser *models.User) error
}
