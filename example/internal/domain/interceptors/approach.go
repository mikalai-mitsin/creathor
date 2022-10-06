package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -destination mock/approach_mock.go github.com/018bf/example/internal/domain/interceptors ApproachInterceptor

type ApproachInterceptor interface {
	Get(ctx context.Context, id string, user *models.User) (*models.Approach, error)
	List(ctx context.Context, filter *models.ApproachFilter, user *models.User) ([]*models.Approach, error)
	Create(ctx context.Context, create *models.ApproachCreate, user *models.User) (*models.Approach, error)
	Update(ctx context.Context, update *models.ApproachUpdate, user *models.User) (*models.Approach, error)
	Delete(ctx context.Context, id string, user *models.User) error
}
