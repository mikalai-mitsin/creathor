package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type UserInterceptor struct {
	userUseCase usecases.UserUseCase
	logger      log.Logger
}

func NewUserInterceptor(
	userUseCase usecases.UserUseCase,
	logger log.Logger,
) interceptors.UserInterceptor {
	return &UserInterceptor{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

func (i *UserInterceptor) Get(ctx context.Context, id string, _ *models.User) (*models.User, error) {
	user, err := i.userUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserInterceptor) List(
	ctx context.Context,
	filter *models.UserFilter,
	_ *models.User,
) ([]*models.User, error) {
	users, err := i.userUseCase.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (i *UserInterceptor) Create(
	ctx context.Context,
	create *models.UserCreate,
	_ *models.User,
) (*models.User, error) {
	user, err := i.userUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserInterceptor) Update(
	ctx context.Context,
	update *models.UserUpdate,
	_ *models.User,
) (*models.User, error) {
	updatedUser, err := i.userUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (i *UserInterceptor) Delete(ctx context.Context, id string, _ *models.User) error {
	if err := i.userUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
