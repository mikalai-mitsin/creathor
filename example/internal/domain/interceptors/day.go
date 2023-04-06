package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// DayInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/day.go . DayInterceptor
type DayInterceptor interface {
	Get(ctx context.Context, id models.UUID, requestUser *models.User) (*models.Day, error)
	List(
		ctx context.Context,
		filter *models.DayFilter,
		requestUser *models.User,
	) ([]*models.Day, uint64, error)
	Update(
		ctx context.Context,
		update *models.DayUpdate,
		requestUser *models.User,
	) (*models.Day, error)
	Create(
		ctx context.Context,
		create *models.DayCreate,
		requestUser *models.User,
	) (*models.Day, error)
	Delete(ctx context.Context, id models.UUID, requestUser *models.User) error
}
