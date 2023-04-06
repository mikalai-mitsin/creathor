package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// ArchInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/arch.go . ArchInterceptor
type ArchInterceptor interface {
	Get(ctx context.Context, id models.UUID, requestUser *models.User) (*models.Arch, error)
	List(
		ctx context.Context,
		filter *models.ArchFilter,
		requestUser *models.User,
	) ([]*models.Arch, uint64, error)
	Update(
		ctx context.Context,
		update *models.ArchUpdate,
		requestUser *models.User,
	) (*models.Arch, error)
	Create(
		ctx context.Context,
		create *models.ArchCreate,
		requestUser *models.User,
	) (*models.Arch, error)
	Delete(ctx context.Context, id models.UUID, requestUser *models.User) error
}
