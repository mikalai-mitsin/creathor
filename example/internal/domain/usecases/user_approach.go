package usecases

import (
	"context"

	"github.com/018bf/creathor/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/user_approach_mock.go github.com/018bf/creathor/internal/domain/usecases UserApproachUseCase

type UserApproachUseCase interface {
	Get(ctx context.Context, id string) (*models.UserApproach, error)
	List(ctx context.Context, filter *models.UserApproachFilter) ([]*models.UserApproach, uint64, error)
	Create(ctx context.Context, create *models.UserApproachCreate) (*models.UserApproach, error)
	Update(ctx context.Context, update *models.UserApproachUpdate) (*models.UserApproach, error)
	Delete(ctx context.Context, id string) error
}
