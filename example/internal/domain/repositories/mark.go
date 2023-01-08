package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/mark_mock.go github.com/018bf/example/internal/domain/repositories MarkRepository

type MarkRepository interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Mark, error)
	List(
		ctx context.Context,
		filter *models.MarkFilter,
	) ([]*models.Mark, error)
	Count(
		ctx context.Context,
		filter *models.MarkFilter,
	) (uint64, error)
	Create(
		ctx context.Context,
		mark *models.Mark,
	) error
	Update(
		ctx context.Context,
		mark *models.Mark,
	) error
	Delete(
		ctx context.Context,
		id string,
	) error
}
