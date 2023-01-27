package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/day_mock.go github.com/018bf/example/internal/domain/usecases DayUseCase

type DayUseCase interface {
	Get(
		ctx context.Context,
		id models.UUID,
	) (*models.Day, error)
	List(
		ctx context.Context,
		filter *models.DayFilter,
	) ([]*models.Day, uint64, error)
	Create(
		ctx context.Context,
		create *models.DayCreate,
	) (*models.Day, error)
	Update(
		ctx context.Context,
		update *models.DayUpdate,
	) (*models.Day, error)
	Delete(
		ctx context.Context,
		id models.UUID,
	) error
}
