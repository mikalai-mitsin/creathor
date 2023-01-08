package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/mark_mock.go github.com/018bf/example/internal/domain/usecases MarkUseCase

type MarkUseCase interface {
	Get(
		ctx context.Context,
		id string,
	) (*models.Mark, error)
	List(
		ctx context.Context,
		filter *models.MarkFilter,
	) ([]*models.Mark, uint64, error)
	Create(
		ctx context.Context,
		create *models.MarkCreate,
	) (*models.Mark, error)
	Update(
		ctx context.Context,
		update *models.MarkUpdate,
	) (*models.Mark, error)
	Delete(
		ctx context.Context,
		id string,
	) error
}
