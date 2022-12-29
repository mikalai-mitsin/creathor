package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -destination mock/session_mock.go github.com/018bf/example/internal/domain/interceptors SessionInterceptor

type SessionInterceptor interface {
	Get(ctx context.Context, id string, user *models.User) (*models.Session, error)
	List(ctx context.Context, filter *models.SessionFilter, user *models.User) ([]*models.Session, uint64, error)
	Create(ctx context.Context, create *models.SessionCreate, user *models.User) (*models.Session, error)
	Update(ctx context.Context, update *models.SessionUpdate, user *models.User) (*models.Session, error)
	Delete(ctx context.Context, id string, user *models.User) error
}
