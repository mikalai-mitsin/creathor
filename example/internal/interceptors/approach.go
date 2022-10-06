package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

//nolint: lll
//go:generate mockgen -destination mock/approach_mock.go github.com/018bf/example/internal/interceptors ApproachInterceptor

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

func (i *ApproachInterceptor) Get(ctx context.Context, id string, user *models.User) (*models.Approach, error) {
	approach, err := i.approachUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return approach, nil
}

func (i *ApproachInterceptor) List(
	ctx context.Context,
	filter *models.ApproachFilter,
	user *models.User,
) ([]*models.Approach, error) {
	approachs, err := i.approachUseCase.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return approachs, nil
}

func (i *ApproachInterceptor) Create(
	ctx context.Context,
	create *models.ApproachCreate,
	user *models.User,
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
	user *models.User,
) (*models.Approach, error) {
	approach, err := i.approachUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return approach, nil
}

func (i *ApproachInterceptor) Delete(ctx context.Context, id string, user *models.User) error {
	if err := i.approachUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
