package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

type DayInterceptor struct {
	dayUseCase  usecases.DayUseCase
	logger      log.Logger
	authUseCase usecases.AuthUseCase
}

func NewDayInterceptor(
	dayUseCase usecases.DayUseCase,
	logger log.Logger,
	authUseCase usecases.AuthUseCase,
) interceptors.DayInterceptor {
	return &DayInterceptor{dayUseCase: dayUseCase, logger: logger, authUseCase: authUseCase}
}

func (i *DayInterceptor) Create(
	ctx context.Context,
	create *models.DayCreate,
	requestUser *models.User,
) (*models.Day, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDDayCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDDayCreate, create); err != nil {
		return nil, err
	}
	day, err := i.dayUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return day, nil
}

func (i *DayInterceptor) Get(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) (*models.Day, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDDayDetail); err != nil {
		return nil, err
	}
	day, err := i.dayUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDDayDetail, day); err != nil {
		return nil, err
	}
	return day, nil
}

func (i *DayInterceptor) List(
	ctx context.Context,
	filter *models.DayFilter,
	requestUser *models.User,
) ([]*models.Day, uint64, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDDayList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDDayList, filter); err != nil {
		return nil, 0, err
	}
	listDays, count, err := i.dayUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return listDays, count, nil
}

func (i *DayInterceptor) Update(
	ctx context.Context,
	update *models.DayUpdate,
	requestUser *models.User,
) (*models.Day, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDDayUpdate); err != nil {
		return nil, err
	}
	day, err := i.dayUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDDayUpdate, day); err != nil {
		return nil, err
	}
	updated, err := i.dayUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (i *DayInterceptor) Delete(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDDayDelete); err != nil {
		return err
	}
	day, err := i.dayUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDDayDelete, day); err != nil {
		return err
	}
	if err := i.dayUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
