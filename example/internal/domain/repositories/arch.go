package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// ArchRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/arch.go github.com/018bf/example/internal/domain/repositories ArchRepository
type ArchRepository interface {
	Get(ctx context.Context, id models.UUID) (*models.Arch, error)
	List(ctx context.Context, filter *models.ArchFilter) ([]*models.Arch, error)
	Count(ctx context.Context, filter *models.ArchFilter) (uint64, error)
	Update(ctx context.Context, update *models.Arch) error
	Create(ctx context.Context, create *models.Arch) error
	Delete(ctx context.Context, id models.UUID) error
}
