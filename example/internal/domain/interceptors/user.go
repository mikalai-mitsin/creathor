package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -destination mock/user_mock.go github.com/018bf/example/internal/domain/interceptors UserInterceptor

type UserInterceptor interface {
	Get(ctx context.Context, id string, user *models.User) (*models.User, error)
	List(ctx context.Context, filter *models.UserFilter, user *models.User) ([]*models.User, error)
	Create(ctx context.Context, create *models.UserCreate, user *models.User) (*models.User, error)
	Update(ctx context.Context, update *models.UserUpdate, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id string, user *models.User) error
}
