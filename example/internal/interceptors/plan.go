package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

type PlanInterceptor struct {
	planUseCase usecases.PlanUseCase
	logger      log.Logger
	authUseCase usecases.AuthUseCase
}

func NewPlanInterceptor(
	planUseCase usecases.PlanUseCase,
	logger log.Logger,
	authUseCase usecases.AuthUseCase,
) interceptors.PlanInterceptor {
	return &PlanInterceptor{planUseCase: planUseCase, logger: logger, authUseCase: authUseCase}
}

func (i *PlanInterceptor) Create(
	ctx context.Context,
	create *models.PlanCreate,
	requestUser *models.User,
) (*models.Plan, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDPlanCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDPlanCreate, create); err != nil {
		return nil, err
	}
	plan, err := i.planUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (i *PlanInterceptor) Get(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) (*models.Plan, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDPlanDetail); err != nil {
		return nil, err
	}
	plan, err := i.planUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDPlanDetail, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (i *PlanInterceptor) List(
	ctx context.Context,
	filter *models.PlanFilter,
	requestUser *models.User,
) ([]*models.Plan, uint64, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDPlanList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDPlanList, filter); err != nil {
		return nil, 0, err
	}
	listPlans, count, err := i.planUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return listPlans, count, nil
}

func (i *PlanInterceptor) Update(
	ctx context.Context,
	update *models.PlanUpdate,
	requestUser *models.User,
) (*models.Plan, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDPlanUpdate); err != nil {
		return nil, err
	}
	plan, err := i.planUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDPlanUpdate, plan); err != nil {
		return nil, err
	}
	updated, err := i.planUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (i *PlanInterceptor) Delete(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDPlanDelete); err != nil {
		return err
	}
	plan, err := i.planUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDPlanDelete, plan); err != nil {
		return err
	}
	if err := i.planUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
