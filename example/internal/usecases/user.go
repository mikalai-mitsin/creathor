package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type UserUseCase struct {
	userRepository repositories.UserRepository
	logger         log.Logger
}

func NewUserUseCase(
	userRepository repositories.UserRepository,
	logger log.Logger,
) usecases.UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (u *UserUseCase) Get(ctx context.Context, id string) (*models.User, error) {
	qr, err := u.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return qr, nil
}

func (u *UserUseCase) List(ctx context.Context, filter *models.UserFilter) ([]*models.User, error) {
	qrs, err := u.userRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return qrs, nil
}

func (u *UserUseCase) Create(ctx context.Context, create *models.UserCreate) (*models.User, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	user := &models.User{
		ID: "",
	}

	if err := u.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) Update(ctx context.Context, update *models.UserUpdate) (*models.User, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	user, err := u.userRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := u.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) Delete(ctx context.Context, id string) error {
	if err := u.userRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
