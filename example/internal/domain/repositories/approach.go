package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/approach_mock.go github.com/018bf/example/internal/domain/repositories ApproachRepository

type ApproachRepository interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Approach, error)
	List(
		ctx context.Context,
		filter *models.ApproachFilter,
	) ([]*models.Approach, error)
	Count(
		ctx context.Context,
		filter *models.ApproachFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		approach *models.Approach,
	) error
	Update(
		ctx context.Context,
		approach *models.Approach,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}
