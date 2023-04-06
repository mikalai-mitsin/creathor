package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// SessionInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/session.go . SessionInterceptor
type SessionInterceptor interface {
	Get(ctx context.Context, id models.UUID, requestUser *models.User) (*models.Session, error)
	List(
		ctx context.Context,
		filter *models.SessionFilter,
		requestUser *models.User,
	) ([]*models.Session, uint64, error)
	Update(
		ctx context.Context,
		update *models.SessionUpdate,
		requestUser *models.User,
	) (*models.Session, error)
	Create(
		ctx context.Context,
		create *models.SessionCreate,
		requestUser *models.User,
	) (*models.Session, error)
	Delete(ctx context.Context, id models.UUID, requestUser *models.User) error
}
