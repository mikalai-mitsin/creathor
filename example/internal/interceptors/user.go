package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type UserInterceptor struct {
	userUseCase usecases.UserUseCase
	authUseCase usecases.AuthUseCase
	logger      log.Logger
}

func NewUserInterceptor(
	userUseCase usecases.UserUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.UserInterceptor {
	return &UserInterceptor{
		userUseCase: userUseCase,
		authUseCase: authUseCase,
		logger:      logger,
	}
}

func (i *UserInterceptor) Get(
	ctx context.Context,
	id string,
	requestUser *models.User,
) (*models.User, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserDetail); err != nil {
		return nil, err
	}
	user, err := i.userUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDUserDetail, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserInterceptor) List(
	ctx context.Context,
	filter *models.UserFilter,
	requestUser *models.User,
) ([]*models.User, uint64, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserList); err != nil {
		return nil, 0, err
	}
	users, count, err := i.userUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (i *UserInterceptor) Create(
	ctx context.Context,
	create *models.UserCreate,
	requestUser *models.User,
) (*models.User, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserCreate); err != nil {
		return nil, err
	}
	user, err := i.userUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserInterceptor) Update(
	ctx context.Context,
	update *models.UserUpdate,
	requestUser *models.User,
) (*models.User, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserUpdate); err != nil {
		return nil, err
	}
	user, err := i.userUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDUserUpdate, user)
	if err != nil {
		return nil, err
	}
	user, err = i.userUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserInterceptor) Delete(
	ctx context.Context,
	id string,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDUserDelete); err != nil {
		return err
	}
	user, err := i.userUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	err = i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDUserDelete, user)
	if err != nil {
		return err
	}
	if err := i.userUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
