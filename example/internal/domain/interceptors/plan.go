package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/plan_mock.go github.com/018bf/example/internal/domain/interceptors PlanInterceptor

type PlanInterceptor interface {
	Get(
		ctx context.Context,
		id models.UUID,
		requestUser *models.User,
	) (*models.Plan, error)
	List(
		ctx context.Context,
		filter *models.PlanFilter,
		requestUser *models.User,
	) ([]*models.Plan, uint64, error)
	Create(
		ctx context.Context,
		create *models.PlanCreate,
		requestUser *models.User,
	) (*models.Plan, error)
	Update(
		ctx context.Context,
		update *models.PlanUpdate,
		requestUser *models.User,
	) (*models.Plan, error)
	Delete(
		ctx context.Context,
		id models.UUID,
		requestUser *models.User,
	) error
}
