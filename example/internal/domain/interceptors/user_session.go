package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_session_mock.go github.com/018bf/example/internal/domain/interceptors UserSessionInterceptor

type UserSessionInterceptor interface {
	Get(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) (*models.UserSession, error)
	List(
		ctx context.Context,
		filter *models.UserSessionFilter,
		requestUser *models.User,
	) ([]*models.UserSession, uint64, error)
	Create(
		ctx context.Context,
		create *models.UserSessionCreate,
		requestUser *models.User,
	) (*models.UserSession, error)
	Update(
		ctx context.Context,
		update *models.UserSessionUpdate,
		requestUser *models.User,
	) (*models.UserSession, error)
	Delete(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) error
}
