package interceptors

import (
	"context"

	"github.com/018bf/creathor/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_approach_mock.go github.com/018bf/creathor/internal/domain/interceptors UserApproachInterceptor

type UserApproachInterceptor interface {
	Get(ctx context.Context, id string, requestUser *models.User) (*models.UserApproach, error)
	List(ctx context.Context, filter *models.UserApproachFilter, requestUser *models.User) ([]*models.UserApproach, uint64, error)
	Create(ctx context.Context, create *models.UserApproachCreate, requestUser *models.User) (*models.UserApproach, error)
	Update(ctx context.Context, update *models.UserApproachUpdate, requestUser *models.User) (*models.UserApproach, error)
	Delete(ctx context.Context, id string, requestUser *models.User) error
}
