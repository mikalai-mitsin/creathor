package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/session_mock.go github.com/018bf/example/internal/domain/repositories SessionRepository

type SessionRepository interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Session, *errs.Error)
	List(
		ctx context.Context,
		filter *models.SessionFilter,
	) ([]*models.Session, *errs.Error)
	Count(
		ctx context.Context,
		filter *models.SessionFilter,
	) (uint64, *errs.Error)
	Create(
		ctx context.Context,
		session *models.Session,
	) *errs.Error
	Update(
		ctx context.Context,
		session *models.Session,
	) *errs.Error
	Delete(
		ctx context.Context,
		id string,
	) *errs.Error
}
