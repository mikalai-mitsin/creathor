package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

type SessionInterceptor struct {
	sessionUseCase usecases.SessionUseCase
	logger         log.Logger
	authUseCase    usecases.AuthUseCase
}

func NewSessionInterceptor(
	sessionUseCase usecases.SessionUseCase,
	logger log.Logger,
	authUseCase usecases.AuthUseCase,
) interceptors.SessionInterceptor {
	return &SessionInterceptor{
		sessionUseCase: sessionUseCase,
		logger:         logger,
		authUseCase:    authUseCase,
	}
}

func (i *SessionInterceptor) Create(
	ctx context.Context,
	create *models.SessionCreate,
	requestUser *models.User,
) (*models.Session, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDSessionCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDSessionCreate, create); err != nil {
		return nil, err
	}
	session, err := i.sessionUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (i *SessionInterceptor) Get(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) (*models.Session, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDSessionDetail); err != nil {
		return nil, err
	}
	session, err := i.sessionUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDSessionDetail, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (i *SessionInterceptor) List(
	ctx context.Context,
	filter *models.SessionFilter,
	requestUser *models.User,
) ([]*models.Session, uint64, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDSessionList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDSessionList, filter); err != nil {
		return nil, 0, err
	}
	listSessions, count, err := i.sessionUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return listSessions, count, nil
}

func (i *SessionInterceptor) Update(
	ctx context.Context,
	update *models.SessionUpdate,
	requestUser *models.User,
) (*models.Session, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDSessionUpdate); err != nil {
		return nil, err
	}
	session, err := i.sessionUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDSessionUpdate, session); err != nil {
		return nil, err
	}
	updated, err := i.sessionUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (i *SessionInterceptor) Delete(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDSessionDelete); err != nil {
		return err
	}
	session, err := i.sessionUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDSessionDelete, session); err != nil {
		return err
	}
	if err := i.sessionUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
