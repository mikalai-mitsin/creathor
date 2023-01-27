package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/day_mock.go github.com/018bf/example/internal/domain/repositories DayRepository

type DayRepository interface {
	Get(
		ctx context.Context,
		id models.UUID,
	) (*models.Day, error)
	List(
		ctx context.Context,
		filter *models.DayFilter,
	) ([]*models.Day, error)
	Count(
		ctx context.Context,
		filter *models.DayFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		day *models.Day,
	) error
	Update(
		ctx context.Context,
		day *models.Day,
	) error
	Delete(
		ctx context.Context,
		id models.UUID,
	) error
}
