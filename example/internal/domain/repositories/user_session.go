package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_session_mock.go github.com/018bf/example/internal/domain/repositories UserSessionRepository

type UserSessionRepository interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.UserSession, error)
	List(
		ctx context.Context,
		filter *models.UserSessionFilter,
	) ([]*models.UserSession, error)
	Count(
		ctx context.Context,
		filter *models.UserSessionFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		userSession *models.UserSession,
	) error
	Update(
		ctx context.Context,
		userSession *models.UserSession,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}
