package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type ApproachUseCase struct {
	approachRepository repositories.ApproachRepository
	logger             log.Logger
}

func NewApproachUseCase(
	approachRepository repositories.ApproachRepository,
	logger log.Logger,
) usecases.ApproachUseCase {
	return &ApproachUseCase{
		approachRepository: approachRepository,
		logger:             logger,
	}
}

func (u *ApproachUseCase) Get(ctx context.Context, id string) (*models.Approach, error) {
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
	approachs, err := u.approachRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.approachRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return approachs, count, nil
}

func (u *ApproachUseCase) Create(ctx context.Context, create *models.ApproachCreate) (*models.Approach, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	approach := &models.Approach{
		ID: "",
	}
	if err := u.approachRepository.Create(ctx, approach); err != nil {
		return nil, err
	}
	return approach, nil
}

func (u *ApproachUseCase) Update(ctx context.Context, update *models.ApproachUpdate) (*models.Approach, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	approach, err := u.approachRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
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
