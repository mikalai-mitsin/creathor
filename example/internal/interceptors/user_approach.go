package interceptors

import (
	"context"

	"github.com/018bf/creathor/internal/domain/interceptors"
	"github.com/018bf/creathor/internal/domain/models"
	"github.com/018bf/creathor/internal/domain/usecases"

	"github.com/018bf/creathor/pkg/log"
)

type UserApproachInterceptor struct {
	userApproachUseCase usecases.UserApproachUseCase
	authUseCase         usecases.AuthUseCase
	logger              log.Logger
}

func NewUserApproachInterceptor(
	userApproachUseCase usecases.UserApproachUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.UserApproachInterceptor {
	return &UserApproachInterceptor{
		userApproachUseCase: userApproachUseCase,
		authUseCase:         authUseCase,
		logger:              logger,
	}
}

func (i *UserApproachInterceptor) Get(
	ctx context.Context,
	id string,
	requestUser *models.User,
) (*models.UserApproach, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserApproachDetail); err != nil {
		return nil, err
	}
	userApproach, err := i.userApproachUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachDetail, userApproach)
	if err != nil {
		return nil, err
	}
	return userApproach, nil
}

func (i *UserApproachInterceptor) List(
	ctx context.Context,
	filter *models.UserApproachFilter,
	requestUser *models.User,
) ([]*models.UserApproach, uint64, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserApproachList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachList, filter); err != nil {
		return nil, 0, err
	}
	userApproaches, count, err := i.userApproachUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return userApproaches, count, nil
}

func (i *UserApproachInterceptor) Create(
	ctx context.Context,
	create *models.UserApproachCreate,
	requestUser *models.User,
) (*models.UserApproach, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserApproachCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachCreate, create); err != nil {
		return nil, err
	}
	userApproach, err := i.userApproachUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return userApproach, nil
}

func (i *UserApproachInterceptor) Update(
	ctx context.Context,
	update *models.UserApproachUpdate,
	requestUser *models.User,
) (*models.UserApproach, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate); err != nil {
		return nil, err
	}
	userApproach, err := i.userApproachUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachUpdate, userApproach); err != nil {
		return nil, err
	}
	updatedUserApproach, err := i.userApproachUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedUserApproach, nil
}

func (i *UserApproachInterceptor) Delete(
	ctx context.Context,
	id string,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserApproachDelete); err != nil {
		return err
	}
	userApproach, err := i.userApproachUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	err = i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDUserApproachDelete, userApproach)
	if err != nil {
		return err
	}
	if err := i.userApproachUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
