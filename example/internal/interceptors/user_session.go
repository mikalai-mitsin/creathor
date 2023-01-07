package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type UserSessionInterceptor struct {
	userSessionUseCase usecases.UserSessionUseCase
	authUseCase        usecases.AuthUseCase
	logger             log.Logger
}

func NewUserSessionInterceptor(
	userSessionUseCase usecases.UserSessionUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.UserSessionInterceptor {
	return &UserSessionInterceptor{
		userSessionUseCase: userSessionUseCase,
		authUseCase:        authUseCase,
		logger:             logger,
	}
}

func (i *UserSessionInterceptor) Get(
	ctx context.Context,
	id string,
	requestUser *models.User,
) (*models.UserSession, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionDetail,
	); err != nil {
		return nil, err
	}
	userSession, err := i.userSessionUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionDetail,
		userSession,
	)
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (i *UserSessionInterceptor) List(
	ctx context.Context,
	filter *models.UserSessionFilter,
	requestUser *models.User,
) ([]*models.UserSession, uint64, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionList,
	); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionList,
		filter,
	); err != nil {
		return nil, 0, err
	}
	userSessions, count, err := i.userSessionUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return userSessions, count, nil
}

func (i *UserSessionInterceptor) Create(
	ctx context.Context,
	create *models.UserSessionCreate,
	requestUser *models.User,
) (*models.UserSession, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionCreate,
	); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionCreate,
		create,
	); err != nil {
		return nil, err
	}
	userSession, err := i.userSessionUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (i *UserSessionInterceptor) Update(
	ctx context.Context,
	update *models.UserSessionUpdate,
	requestUser *models.User,
) (*models.UserSession, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionUpdate,
	); err != nil {
		return nil, err
	}
	userSession, err := i.userSessionUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionUpdate,
		userSession,
	); err != nil {
		return nil, err
	}
	updatedUserSession, err := i.userSessionUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedUserSession, nil
}

func (i *UserSessionInterceptor) Delete(
	ctx context.Context,
	id string,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionDelete,
	); err != nil {
		return err
	}
	userSession, err := i.userSessionUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDUserSessionDelete,
		userSession,
	)
	if err != nil {
		return err
	}
	if err := i.userSessionUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
