package usecases

import (
	"context"
	"strings"
	"time"

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
	return &UserUseCase{userRepository: userRepository, clock: clock, logger: logger}
}
func (u *UserUseCase) Get(ctx context.Context, id models.UUID) (*models.User, error) {
	user, err := u.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserUseCase) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.userRepository.GetByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) List(
	ctx context.Context,
	filter *models.UserFilter,
) ([]*models.User, uint64, error) {
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
func (u *UserUseCase) Create(ctx context.Context, create *models.UserCreate) (*models.User, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	user := &models.User{
		ID:        "",
		FirstName: "",
		LastName:  "",
		Password:  "",
		Email:     strings.ToLower(create.Email),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		GroupID:   models.GroupIDUser,
	}
	user.SetPassword(create.Password)
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
	if update.FirstName != nil {
		user.FirstName = *update.FirstName
	}
	if update.LastName != nil {
		user.LastName = *update.LastName
	}
	if update.Password != nil {
		user.SetPassword(*update.Password)
	}
	if update.Email != nil {
		user.Email = strings.ToLower(*update.Email)
	}
	if err := u.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserUseCase) Delete(ctx context.Context, id models.UUID) error {
	if err := u.userRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
