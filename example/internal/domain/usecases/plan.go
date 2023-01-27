package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/plan_mock.go github.com/018bf/example/internal/domain/usecases PlanUseCase

type PlanUseCase interface {
	Get(
		ctx context.Context,
		id models.UUID,
	) (*models.Plan, error)
	List(
		ctx context.Context,
		filter *models.PlanFilter,
	) ([]*models.Plan, uint64, error)
	Create(
		ctx context.Context,
		create *models.PlanCreate,
	) (*models.Plan, error)
	Update(
		ctx context.Context,
		update *models.PlanUpdate,
	) (*models.Plan, error)
	Delete(
		ctx context.Context,
		id models.UUID,
	) error
}
