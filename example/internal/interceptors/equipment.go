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
	authUseCase      usecases.AuthUseCase
	logger           log.Logger
}

func NewEquipmentInterceptor(
	equipmentUseCase usecases.EquipmentUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.EquipmentInterceptor {
	return &EquipmentInterceptor{
		equipmentUseCase: equipmentUseCase,
		authUseCase:      authUseCase,
		logger:           logger,
	}
}

func (i *EquipmentInterceptor) Get(
	ctx context.Context,
	id string,
	requestUser *models.User,
) (*models.Equipment, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentDetail,
	); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentDetail,
		equipment,
	)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
	requestUser *models.User,
) ([]*models.Equipment, uint64, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentList,
	); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentList,
		filter,
	); err != nil {
		return nil, 0, err
	}
	equipment, count, err := i.equipmentUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return equipment, count, nil
}

func (i *EquipmentInterceptor) Create(
	ctx context.Context,
	create *models.EquipmentCreate,
	requestUser *models.User,
) (*models.Equipment, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentCreate,
	); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentCreate,
		create,
	); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (i *EquipmentInterceptor) Update(
	ctx context.Context,
	update *models.EquipmentUpdate,
	requestUser *models.User,
) (*models.Equipment, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentUpdate,
	); err != nil {
		return nil, err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentUpdate,
		equipment,
	); err != nil {
		return nil, err
	}
	updatedEquipment, err := i.equipmentUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedEquipment, nil
}

func (i *EquipmentInterceptor) Delete(
	ctx context.Context,
	id string,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentDelete,
	); err != nil {
		return err
	}
	equipment, err := i.equipmentUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDEquipmentDelete,
		equipment,
	)
	if err != nil {
		return err
	}
	if err := i.equipmentUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
