package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_mock.go github.com/018bf/example/internal/domain/repositories UserRepository

type UserRepository interface {
	Get(
		ctx context.Context,
		id models.UUID,
	) (*models.User, error)
	GetByEmail(
		ctx context.Context,
		email string,
	) (*models.User, error)
	List(
		ctx context.Context,
		filter *models.UserFilter,
	) ([]*models.User, error)
	Count(
		ctx context.Context,
		filter *models.UserFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		user *models.User,
	) error
	Update(
		ctx context.Context,
		user *models.User,
	) error
	Delete(
		ctx context.Context,
		id models.UUID,
	) error
}
