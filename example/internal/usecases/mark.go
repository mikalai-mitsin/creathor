package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type MarkUseCase struct {
	markRepository repositories.MarkRepository
	clock          clock.Clock
	logger         log.Logger
}

func NewMarkUseCase(
	markRepository repositories.MarkRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.MarkUseCase {
	return &MarkUseCase{
		markRepository: markRepository,
		clock:          clock,
		logger:         logger,
	}
}

func (u *MarkUseCase) Get(
	ctx context.Context,
	id string,
) (*models.Mark, error) {
	mark, err := u.markRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return mark, nil
}

func (u *MarkUseCase) List(
	ctx context.Context,
	filter *models.MarkFilter,
) ([]*models.Mark, uint64, error) {
	marks, err := u.markRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.markRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return marks, count, nil
}

func (u *MarkUseCase) Create(
	ctx context.Context,
	create *models.MarkCreate,
) (*models.Mark, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	mark := &models.Mark{
		ID:        "",
		Name:      create.Name,
		Title:     create.Title,
		Weight:    create.Weight,
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := u.markRepository.Create(ctx, mark); err != nil {
		return nil, err
	}
	return mark, nil
}

func (u *MarkUseCase) Update(
	ctx context.Context,
	update *models.MarkUpdate,
) (*models.Mark, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	mark, err := u.markRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if update.Name != nil {
		mark.Name = *update.Name
	}
	if update.Title != nil {
		mark.Title = *update.Title
	}
	if update.Weight != nil {
		mark.Weight = *update.Weight
	}
	mark.UpdatedAt = u.clock.Now()
	if err := u.markRepository.Update(ctx, mark); err != nil {
		return nil, err
	}
	return mark, nil
}

func (u *MarkUseCase) Delete(ctx context.Context, id string) error {
	if err := u.markRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
