package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/plan_mock.go github.com/018bf/example/internal/domain/repositories PlanRepository

type PlanRepository interface {
	Get(
		ctx context.Context,
		id models.UUID,
	) (*models.Plan, error)
	List(
		ctx context.Context,
		filter *models.PlanFilter,
	) ([]*models.Plan, error)
	Count(
		ctx context.Context,
		filter *models.PlanFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		plan *models.Plan,
	) error
	Update(
		ctx context.Context,
		plan *models.Plan,
	) error
	Delete(
		ctx context.Context,
		id models.UUID,
	) error
}
