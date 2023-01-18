package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/session_mock.go github.com/018bf/example/internal/domain/repositories SessionRepository

type SessionRepository interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Session, error)
	List(
		ctx context.Context,
		filter *models.SessionFilter,
	) ([]*models.Session, error)
	Count(
		ctx context.Context,
		filter *models.SessionFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		session *models.Session,
	) error
	Update(
		ctx context.Context,
		session *models.Session,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}
