package usecases

import (
	"context"

	"github.com/018bf/creathor/internal/domain/models"
	"github.com/018bf/creathor/internal/domain/repositories"
	"github.com/018bf/creathor/internal/domain/usecases"

	"github.com/018bf/creathor/pkg/log"
)

type UserApproachUseCase struct {
	userApproachRepository repositories.UserApproachRepository
	logger                 log.Logger
}

func NewUserApproachUseCase(
	userApproachRepository repositories.UserApproachRepository,
	logger log.Logger,
) usecases.UserApproachUseCase {
	return &UserApproachUseCase{
		userApproachRepository: userApproachRepository,
		logger:                 logger,
	}
}

func (u *UserApproachUseCase) Get(ctx context.Context, id string) (*models.UserApproach, error) {
	userApproach, err := u.userApproachRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return userApproach, nil
}

func (u *UserApproachUseCase) List(
	ctx context.Context,
	filter *models.UserApproachFilter,
) ([]*models.UserApproach, uint64, error) {
	userApproaches, err := u.userApproachRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.userApproachRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return userApproaches, count, nil
}

func (u *UserApproachUseCase) Create(ctx context.Context, create *models.UserApproachCreate) (*models.UserApproach, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	userApproach := &models.UserApproach{
		ID: "",
	}
	if err := u.userApproachRepository.Create(ctx, userApproach); err != nil {
		return nil, err
	}
	return userApproach, nil
}

func (u *UserApproachUseCase) Update(ctx context.Context, update *models.UserApproachUpdate) (*models.UserApproach, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	userApproach, err := u.userApproachRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := u.userApproachRepository.Update(ctx, userApproach); err != nil {
		return nil, err
	}
	return userApproach, nil
}

func (u *UserApproachUseCase) Delete(ctx context.Context, id string) error {
	if err := u.userApproachRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
