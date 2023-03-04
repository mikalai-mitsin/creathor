package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// EquipmentUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/equipment.go github.com/018bf/example/internal/domain/usecases EquipmentUseCase
type EquipmentUseCase interface {
	Get(ctx context.Context, id models.UUID) (*models.Equipment, error)
	List(ctx context.Context, filter *models.EquipmentFilter) ([]*models.Equipment, uint64, error)
	Update(ctx context.Context, update *models.EquipmentUpdate) (*models.Equipment, error)
	Create(ctx context.Context, create *models.EquipmentCreate) (*models.Equipment, error)
	Delete(ctx context.Context, id models.UUID) error
}
