package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -destination mock/equipment_mock.go github.com/018bf/example/internal/domain/repositories EquipmentRepository

type EquipmentRepository interface {
	Get(ctx context.Context, id string) (*models.Equipment, error)
	List(ctx context.Context, filter *models.EquipmentFilter) ([]*models.Equipment, error)
	Create(ctx context.Context, equipment *models.Equipment) error
	Update(ctx context.Context, equipment *models.Equipment) error
	Delete(ctx context.Context, id string) error
}
