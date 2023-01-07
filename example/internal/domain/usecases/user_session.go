package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_session_mock.go github.com/018bf/example/internal/domain/usecases UserSessionUseCase

type UserSessionUseCase interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.UserSession, error)
	List(
		ctx context.Context,
		filter *models.UserSessionFilter,
	) ([]*models.UserSession, uint64, error)
	Create(
		ctx context.Context,
		create *models.UserSessionCreate,
	) (*models.UserSession, error)
	Update(
		ctx context.Context,
		update *models.UserSessionUpdate,
	) (*models.UserSession, error)
	Delete(
		ctx context.Context,
		id string,
	) error
}
