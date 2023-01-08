package interceptors

import (
	"context"
	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type AuthInterceptor struct {
	authUseCase usecases.AuthUseCase
	userUseCase usecases.UserUseCase
	clock       clock.Clock
	logger      log.Logger
}

func NewAuthInterceptor(
	authUseCase usecases.AuthUseCase,
	userUseCase usecases.UserUseCase,
	clock clock.Clock,
	logger log.Logger,
) interceptors.AuthInterceptor {
	return &AuthInterceptor{
		authUseCase: authUseCase,
		userUseCase: userUseCase,
		clock:       clock,
		logger:      logger,
	}
}

func (i *AuthInterceptor) CreateToken(
	ctx context.Context,
	login *models.Login,
) (*models.TokenPair, error) {
	pair, err := i.authUseCase.CreateToken(ctx, login)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (i *AuthInterceptor) ValidateToken(
	ctx context.Context,
	token models.Token,
) error {
	if err := i.authUseCase.ValidateToken(ctx, token); err != nil {
		return err
	}
	return nil
}

func (i *AuthInterceptor) RefreshToken(
	ctx context.Context,
	refresh models.Token,
) (*models.TokenPair, error) {
	pair, err := i.authUseCase.RefreshToken(ctx, refresh)
	if err != nil {
		return nil, err
	}
	return pair, nil
}

func (i *AuthInterceptor) Auth(
	ctx context.Context,
	access models.Token,
) (*models.User, error) {
	user, err := i.authUseCase.Auth(ctx, access)
	if err != nil {
		return nil, err
	}
	return user, nil
}
