package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/arch_mock.go github.com/018bf/example/internal/domain/interceptors ArchInterceptor

type ArchInterceptor interface {
	Get(
		ctx context.Context,
		id models.UUID,
		requestUser *models.User,
	) (*models.Arch, error)
	List(
		ctx context.Context,
		filter *models.ArchFilter,
		requestUser *models.User,
	) ([]*models.Arch, uint64, error)
	Create(
		ctx context.Context,
		create *models.ArchCreate,
		requestUser *models.User,
	) (*models.Arch, error)
	Update(
		ctx context.Context,
		update *models.ArchUpdate,
		requestUser *models.User,
	) (*models.Arch, error)
	Delete(
		ctx context.Context,
		id models.UUID,
		requestUser *models.User,
	) error
}
