package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// SessionUseCase - domain layer use case interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/session.go github.com/018bf/example/internal/domain/usecases SessionUseCase
type SessionUseCase interface {
	Get(ctx context.Context, id models.UUID) (*models.Session, error)
	List(ctx context.Context, filter *models.SessionFilter) ([]*models.Session, uint64, error)
	Update(ctx context.Context, update *models.SessionUpdate) (*models.Session, error)
	Create(ctx context.Context, create *models.SessionCreate) (*models.Session, error)
	Delete(ctx context.Context, id models.UUID) error
}
