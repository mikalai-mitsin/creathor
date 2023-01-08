package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_mock.go github.com/018bf/example/internal/domain/interceptors UserInterceptor

type UserInterceptor interface {
	Get(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) (*models.User, error)
	List(
		ctx context.Context,
		filter *models.UserFilter,
		requestUser *models.User,
	) ([]*models.User, uint64, error)
	Create(
		ctx context.Context,
		create *models.UserCreate,
		requestUser *models.User,
	) (*models.User, error)
	Update(
		ctx context.Context,
		update *models.UserUpdate,
		requestUser *models.User,
	) (*models.User, error)
	Delete(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) error
}
