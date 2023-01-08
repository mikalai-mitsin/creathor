package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type SessionUseCase struct {
	sessionRepository repositories.SessionRepository
	clock             clock.Clock
	logger            log.Logger
}

func NewSessionUseCase(
	sessionRepository repositories.SessionRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.SessionUseCase {
	return &SessionUseCase{
		sessionRepository: sessionRepository,
		clock:             clock,
		logger:            logger,
	}
}

func (u *SessionUseCase) Get(
	ctx context.Context,
	id string,
) (*models.Session, error) {
	session, err := u.sessionRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (u *SessionUseCase) List(
	ctx context.Context,
	filter *models.SessionFilter,
) ([]*models.Session, uint64, error) {
	sessions, err := u.sessionRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.sessionRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return sessions, count, nil
}

func (u *SessionUseCase) Create(
	ctx context.Context,
	create *models.SessionCreate,
) (*models.Session, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	session := &models.Session{
		ID:        "",
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := u.sessionRepository.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (u *SessionUseCase) Update(
	ctx context.Context,
	update *models.SessionUpdate,
) (*models.Session, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	session, err := u.sessionRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	session.UpdatedAt = u.clock.Now()
	if err := u.sessionRepository.Update(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (u *SessionUseCase) Delete(ctx context.Context, id string) error {
	if err := u.sessionRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
