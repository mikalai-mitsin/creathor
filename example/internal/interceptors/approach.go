package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type ApproachInterceptor struct {
	approachUseCase usecases.ApproachUseCase
	logger          log.Logger
}

func NewApproachInterceptor(
	approachUseCase usecases.ApproachUseCase,
	logger log.Logger,
) interceptors.ApproachInterceptor {
	return &ApproachInterceptor{
		approachUseCase: approachUseCase,
		logger:          logger,
	}
}

func (i *ApproachInterceptor) Get(
	ctx context.Context,
	id string,
) (*models.Approach, error) {
	approach, err := i.approachUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return approach, nil
}

func (i *ApproachInterceptor) List(
	ctx context.Context,
	filter *models.ApproachFilter,
) ([]*models.Approach, uint64, error) {
	approachs, count, err := i.approachUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return approachs, count, nil
}

func (i *ApproachInterceptor) Create(
	ctx context.Context,
	create *models.ApproachCreate,
) (*models.Approach, error) {
	approach, err := i.approachUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return approach, nil
}

func (i *ApproachInterceptor) Update(
	ctx context.Context,
	update *models.ApproachUpdate,
) (*models.Approach, error) {
	updatedApproach, err := i.approachUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedApproach, nil
}

func (i *ApproachInterceptor) Delete(
	ctx context.Context,
	id string,
) error {
	if err := i.approachUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
