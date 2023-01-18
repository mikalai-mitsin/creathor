package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type UserUseCase struct {
	userRepository repositories.UserRepository
	clock          clock.Clock
	logger         log.Logger
}

func NewUserUseCase(
	userRepository repositories.UserRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		clock:          clock,
		logger:         logger,
	}
}

func (u *UserUseCase) Get(
	ctx context.Context,
	id string,
) (*models.User, *errs.Error) {
	user, err := u.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) List(
	ctx context.Context,
	filter *models.UserFilter,
) ([]*models.User, uint64, *errs.Error) {
	users, err := u.userRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.userRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u *UserUseCase) Create(
	ctx context.Context,
	create *models.UserCreate,
) (*models.User, *errs.Error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	user := &models.User{
		ID:        "",
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := u.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) Update(
	ctx context.Context,
	update *models.UserUpdate,
) (*models.User, *errs.Error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	user, err := u.userRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	user.UpdatedAt = u.clock.Now()
	if err := u.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) Delete(ctx context.Context, id string) *errs.Error {
	if err := u.userRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
