package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

type AuthUseCase struct {
	authRepository       repositories.AuthRepository
	userRepository       repositories.UserRepository
	permissionRepository repositories.PermissionRepository
	logger               log.Logger
}

func NewAuthUseCase(
	authRepository repositories.AuthRepository,
	userRepository repositories.UserRepository,
	permissionRepository repositories.PermissionRepository,
	logger log.Logger,
) usecases.AuthUseCase {
	return &AuthUseCase{
		authRepository:       authRepository,
		userRepository:       userRepository,
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (u AuthUseCase) CreateToken(
	ctx context.Context,
	login *models.Login,
) (*models.TokenPair, error) {
	user, err := u.userRepository.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, err
	}
	if err := user.CheckPassword(login.Password); err != nil {
		return nil, err
	}
	tokenPair, err := u.authRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return tokenPair, nil
}

func (u AuthUseCase) CreateTokenByUser(
	ctx context.Context,
	user *models.User,
) (*models.TokenPair, error) {
	tokenPair, err := u.authRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return tokenPair, nil
}

func (u AuthUseCase) RefreshToken(
	ctx context.Context,
	refresh models.Token,
) (*models.TokenPair, error) {
	pair, err := u.authRepository.RefreshToken(ctx, refresh)
	if err != nil {
		return nil, err
	}
	return pair, nil
}
func (u AuthUseCase) ValidateToken(ctx context.Context, access models.Token) error {
	if err := u.authRepository.Validate(ctx, access); err != nil {
		return err
	}
	return nil
}
func (u AuthUseCase) Auth(ctx context.Context, access models.Token) (*models.User, error) {
	userID, err := u.authRepository.GetSubject(ctx, access)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepository.Get(ctx, models.UUID(userID))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u AuthUseCase) HasPermission(
	ctx context.Context,
	user *models.User,
	permission models.PermissionID,
) error {
	if err := u.permissionRepository.HasPermission(ctx, permission, user); err != nil {
		return err
	}
	return nil
}

func (u AuthUseCase) HasObjectPermission(
	ctx context.Context,
	user *models.User,
	permission models.PermissionID,
	object any,
) error {
	if err := u.permissionRepository.HasObjectPermission(ctx, permission, user, object); err != nil {
		return err
	}
	return nil
}
