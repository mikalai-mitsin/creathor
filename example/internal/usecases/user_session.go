package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type UserSessionUseCase struct {
	userSessionRepository repositories.UserSessionRepository
	clock                 clock.Clock
	logger                log.Logger
}

func NewUserSessionUseCase(
	userSessionRepository repositories.UserSessionRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.UserSessionUseCase {
	return &UserSessionUseCase{
		userSessionRepository: userSessionRepository,
		clock:                 clock,
		logger:                logger,
	}
}

func (u *UserSessionUseCase) Get(
	ctx context.Context,
	id string,
) (*models.UserSession, error) {
	userSession, err := u.userSessionRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (u *UserSessionUseCase) List(
	ctx context.Context,
	filter *models.UserSessionFilter,
) ([]*models.UserSession, uint64, error) {
	userSessions, err := u.userSessionRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.userSessionRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return userSessions, count, nil
}

func (u *UserSessionUseCase) Create(
	ctx context.Context,
	create *models.UserSessionCreate,
) (*models.UserSession, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	userSession := &models.UserSession{
		ID:        "",
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := u.userSessionRepository.Create(ctx, userSession); err != nil {
		return nil, err
	}
	return userSession, nil
}

func (u *UserSessionUseCase) Update(
	ctx context.Context,
	update *models.UserSessionUpdate,
) (*models.UserSession, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	userSession, err := u.userSessionRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	userSession.UpdatedAt = u.clock.Now()
	if err := u.userSessionRepository.Update(ctx, userSession); err != nil {
		return nil, err
	}
	return userSession, nil
}

func (u *UserSessionUseCase) Delete(ctx context.Context, id string) error {
	if err := u.userSessionRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
