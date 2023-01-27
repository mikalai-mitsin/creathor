package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/day_mock.go github.com/018bf/example/internal/domain/interceptors DayInterceptor

type DayInterceptor interface {
	Get(
		ctx context.Context,
		id models.UUID,
		requestUser *models.User,
	) (*models.Day, error)
	List(
		ctx context.Context,
		filter *models.DayFilter,
		requestUser *models.User,
	) ([]*models.Day, uint64, error)
	Create(
		ctx context.Context,
		create *models.DayCreate,
		requestUser *models.User,
	) (*models.Day, error)
	Update(
		ctx context.Context,
		update *models.DayUpdate,
		requestUser *models.User,
	) (*models.Day, error)
	Delete(
		ctx context.Context,
		id models.UUID,
		requestUser *models.User,
	) error
}
