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
}

func NewSessionInterceptor(
	sessionUseCase usecases.SessionUseCase,
	logger log.Logger,
) interceptors.SessionInterceptor {
	return &SessionInterceptor{
		sessionUseCase: sessionUseCase,
		logger:         logger,
	}
}

func (i *SessionInterceptor) Get(
	ctx context.Context,
	id string,
) (*models.Session, error) {
	session, err := i.sessionUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (i *SessionInterceptor) List(
	ctx context.Context,
	filter *models.SessionFilter,
) ([]*models.Session, uint64, error) {
	sessions, count, err := i.sessionUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return sessions, count, nil
}

func (i *SessionInterceptor) Create(
	ctx context.Context,
	create *models.SessionCreate,
) (*models.Session, error) {
	session, err := i.sessionUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (i *SessionInterceptor) Update(
	ctx context.Context,
	update *models.SessionUpdate,
) (*models.Session, error) {
	updatedSession, err := i.sessionUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedSession, nil
}

func (i *SessionInterceptor) Delete(
	ctx context.Context,
	id string,
) error {
	if err := i.sessionUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
