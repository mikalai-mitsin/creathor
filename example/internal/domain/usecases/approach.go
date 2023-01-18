package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/approach_mock.go github.com/018bf/example/internal/domain/usecases ApproachUseCase

type ApproachUseCase interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Approach, *errs.Error)
	List(
		ctx context.Context,
		filter *models.ApproachFilter,
	) ([]*models.Approach, uint64, *errs.Error)
	Create(
		ctx context.Context,
		create *models.ApproachCreate,
	) (*models.Approach, *errs.Error)
	Update(
		ctx context.Context,
		update *models.ApproachUpdate,
	) (*models.Approach, *errs.Error)
	Delete(
		ctx context.Context,
		id string,
	) *errs.Error
}
