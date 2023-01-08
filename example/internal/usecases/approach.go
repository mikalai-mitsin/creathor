package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type ApproachUseCase struct {
	approachRepository repositories.ApproachRepository
	clock              clock.Clock
	logger             log.Logger
}

func NewApproachUseCase(
	approachRepository repositories.ApproachRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.ApproachUseCase {
	return &ApproachUseCase{
		approachRepository: approachRepository,
		clock:              clock,
		logger:             logger,
	}
}

func (u *ApproachUseCase) Get(
	ctx context.Context,
	id string,
) (*models.Approach, error) {
	approach, err := u.approachRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return approach, nil
}

func (u *ApproachUseCase) List(
	ctx context.Context,
	filter *models.ApproachFilter,
) ([]*models.Approach, uint64, error) {
	approaches, err := u.approachRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.approachRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return approaches, count, nil
}

func (u *ApproachUseCase) Create(
	ctx context.Context,
	create *models.ApproachCreate,
) (*models.Approach, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	approach := &models.Approach{
		ID:        "",
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := u.approachRepository.Create(ctx, approach); err != nil {
		return nil, err
	}
	return approach, nil
}

func (u *ApproachUseCase) Update(
	ctx context.Context,
	update *models.ApproachUpdate,
) (*models.Approach, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	approach, err := u.approachRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	approach.UpdatedAt = u.clock.Now()
	if err := u.approachRepository.Update(ctx, approach); err != nil {
		return nil, err
	}
	return approach, nil
}

func (u *ApproachUseCase) Delete(ctx context.Context, id string) error {
	if err := u.approachRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
