package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// SessionRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/session.go github.com/018bf/example/internal/domain/repositories SessionRepository
type SessionRepository interface {
	Get(ctx context.Context, id models.UUID) (*models.Session, error)
	List(ctx context.Context, filter *models.SessionFilter) ([]*models.Session, error)
	Count(ctx context.Context, filter *models.SessionFilter) (uint64, error)
	Update(ctx context.Context, update *models.Session) error
	Create(ctx context.Context, create *models.Session) error
	Delete(ctx context.Context, id models.UUID) error
}
