package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/equipment_mock.go github.com/018bf/example/internal/domain/usecases EquipmentUseCase

type EquipmentUseCase interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Equipment, *errs.Error)
	List(
		ctx context.Context,
		filter *models.EquipmentFilter,
	) ([]*models.Equipment, uint64, *errs.Error)
	Create(
		ctx context.Context,
		create *models.EquipmentCreate,
	) (*models.Equipment, *errs.Error)
	Update(
		ctx context.Context,
		update *models.EquipmentUpdate,
	) (*models.Equipment, *errs.Error)
	Delete(
		ctx context.Context,
		id string,
	) *errs.Error
}
