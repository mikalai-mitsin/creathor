package interceptors

import (
	"context"

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

func (i *EquipmentInterceptor) Get(ctx context.Context, id string, _ *models.User) (*models.Equipment, error) {
	equipment, err := i.equipmentUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
	_ *models.User,
) ([]*models.Equipment, error) {
	equipments, err := i.equipmentUseCase.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return equipments, nil
}

func (i *EquipmentInterceptor) Create(
	ctx context.Context,
	create *models.EquipmentCreate,
	_ *models.User,
) (*models.Equipment, error) {
	equipment, err := i.equipmentUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) Update(
	ctx context.Context,
	update *models.EquipmentUpdate,
	_ *models.User,
) (*models.Equipment, error) {
	updatedEquipment, err := i.equipmentUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedEquipment, nil
}

func (i *EquipmentInterceptor) Delete(ctx context.Context, id string, _ *models.User) error {
	if err := i.equipmentUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
