package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/equipment_mock.go github.com/018bf/example/internal/domain/usecases EquipmentUseCase

type EquipmentUseCase interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Equipment, error)
	List(
		ctx context.Context,
		filter *models.EquipmentFilter,
	) ([]*models.Equipment, uint64, error)
	Create(
		ctx context.Context,
		create *models.EquipmentCreate,
	) (*models.Equipment, error)
	Update(
		ctx context.Context,
		update *models.EquipmentUpdate,
	) (*models.Equipment, error)
	Delete(
		ctx context.Context,
		id string,
	) error
}
