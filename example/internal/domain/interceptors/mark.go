package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/mark_mock.go github.com/018bf/example/internal/domain/interceptors MarkInterceptor

type MarkInterceptor interface {
	Get(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) (*models.Mark, error)
	List(
		ctx context.Context,
		filter *models.MarkFilter,
		requestUser *models.User,
	) ([]*models.Mark, uint64, error)
	Create(
		ctx context.Context,
		create *models.MarkCreate,
		requestUser *models.User,
	) (*models.Mark, error)
	Update(
		ctx context.Context,
		update *models.MarkUpdate,
		requestUser *models.User,
	) (*models.Mark, error)
	Delete(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) error
}
