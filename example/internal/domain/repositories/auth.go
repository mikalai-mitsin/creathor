package repositories

import (
	"context"
	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -build_flags=-mod=mod -destination mock/auth_mock.go github.com/018bf/example/internal/domain/repositories AuthRepository

type AuthRepository interface {
	Create(
		ctx context.Context,
		user *models.User,
	) (*models.TokenPair, error)
	Validate(
		ctx context.Context,
		token models.Token,
	) error
	RefreshToken(
		ctx context.Context,
		token models.Token,
	) (*models.TokenPair, error)
	GetSubject(
		ctx context.Context,
		token models.Token,
	) (string, error)
}
