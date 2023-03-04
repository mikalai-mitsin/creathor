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
	authUseCase      usecases.AuthUseCase
}

func NewEquipmentInterceptor(
	equipmentUseCase usecases.EquipmentUseCase,
	logger log.Logger,
	authUseCase usecases.AuthUseCase,
) interceptors.EquipmentInterceptor {
	return &EquipmentInterceptor{
		equipmentUseCase: equipmentUseCase,
		logger:           logger,
		authUseCase:      authUseCase,
	}
}

func (i *EquipmentInterceptor) Create(
	ctx context.Context,
	create *models.EquipmentCreate,
	requestUser *models.User,
) (*models.Equipment, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDEquipmentCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentCreate, create); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) Get(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) (*models.Equipment, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDEquipmentDetail); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentDetail, equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
	requestUser *models.User,
) ([]*models.Equipment, uint64, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDEquipmentList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentList, filter); err != nil {
		return nil, 0, err
	}
	listEquipment, count, err := i.equipmentUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return listEquipment, count, nil
}

func (i *EquipmentInterceptor) Update(
	ctx context.Context,
	update *models.EquipmentUpdate,
	requestUser *models.User,
) (*models.Equipment, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentUpdate, equipment); err != nil {
		return nil, err
	}
	updated, err := i.equipmentUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (i *EquipmentInterceptor) Delete(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDEquipmentDelete); err != nil {
		return err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDEquipmentDelete, equipment); err != nil {
		return err
	}
	if err := i.equipmentUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
