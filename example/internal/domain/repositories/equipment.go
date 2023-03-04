package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// EquipmentRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/equipment.go github.com/018bf/example/internal/domain/repositories EquipmentRepository
type EquipmentRepository interface {
	Get(ctx context.Context, id models.UUID) (*models.Equipment, error)
	List(ctx context.Context, filter *models.EquipmentFilter) ([]*models.Equipment, error)
	Count(ctx context.Context, filter *models.EquipmentFilter) (uint64, error)
	Update(ctx context.Context, update *models.Equipment) error
	Create(ctx context.Context, create *models.Equipment) error
	Delete(ctx context.Context, id models.UUID) error
}
