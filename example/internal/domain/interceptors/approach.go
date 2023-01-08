package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/approach_mock.go github.com/018bf/example/internal/domain/interceptors ApproachInterceptor

type ApproachInterceptor interface {
	Get(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) (*models.Approach, error)
	List(
		ctx context.Context,
		filter *models.ApproachFilter,
		requestUser *models.User,
	) ([]*models.Approach, uint64, error)
	Create(
		ctx context.Context,
		create *models.ApproachCreate,
		requestUser *models.User,
	) (*models.Approach, error)
	Update(
		ctx context.Context,
		update *models.ApproachUpdate,
		requestUser *models.User,
	) (*models.Approach, error)
	Delete(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) error
}
