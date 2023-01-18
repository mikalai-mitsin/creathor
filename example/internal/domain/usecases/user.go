package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_mock.go github.com/018bf/example/internal/domain/usecases UserUseCase

type UserUseCase interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.User, error)
	GetByEmail(
		ctx context.Context,
		email string,
	) (*models.User, error)
	List(
		ctx context.Context,
		filter *models.UserFilter,
	) ([]*models.User, uint64, error)
	Create(
		ctx context.Context,
		create *models.UserCreate,
	) (*models.User, error)
	Update(
		ctx context.Context,
		update *models.UserUpdate,
	) (*models.User, error)
	Delete(
		ctx context.Context,
		id string,
	) error
}
