package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_mock.go github.com/018bf/example/internal/domain/repositories UserRepository

type UserRepository interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.User, *errs.Error)
	List(
		ctx context.Context,
		filter *models.UserFilter,
	) ([]*models.User, *errs.Error)
	Count(
		ctx context.Context,
		filter *models.UserFilter,
	) (uint64, *errs.Error)
	Create(
		ctx context.Context,
		user *models.User,
	) *errs.Error
	Update(
		ctx context.Context,
		user *models.User,
	) *errs.Error
	Delete(
		ctx context.Context,
		id string,
	) *errs.Error
}
