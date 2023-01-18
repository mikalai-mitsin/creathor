package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type EquipmentInterceptor struct {
	equipmentUseCase usecases.EquipmentUseCase
	logger           log.Logger
}

func NewEquipmentInterceptor(
	equipmentUseCase usecases.EquipmentUseCase,
	logger log.Logger,
) interceptors.EquipmentInterceptor {
	return &EquipmentInterceptor{
		equipmentUseCase: equipmentUseCase,
		logger:           logger,
	}
}

func (i *EquipmentInterceptor) Get(
	ctx context.Context,
	id string,
) (*models.Equipment, *errs.Error) {
	equipment, err := i.equipmentUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
) ([]*models.Equipment, uint64, *errs.Error) {
	equipment, count, err := i.equipmentUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return equipment, count, nil
}

func (i *EquipmentInterceptor) Create(
	ctx context.Context,
	create *models.EquipmentCreate,
) (*models.Equipment, *errs.Error) {
	equipment, err := i.equipmentUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) Update(
	ctx context.Context,
	update *models.EquipmentUpdate,
) (*models.Equipment, *errs.Error) {
	updatedEquipment, err := i.equipmentUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedEquipment, nil
}

func (i *EquipmentInterceptor) Delete(
	ctx context.Context,
	id string,
) *errs.Error {
	if err := i.equipmentUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
