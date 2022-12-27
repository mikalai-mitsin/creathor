package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -destination mock/equipment_mock.go github.com/018bf/example/internal/domain/interceptors EquipmentInterceptor

type EquipmentInterceptor interface {
	Get(ctx context.Context, id string, user *models.User) (*models.Equipment, error)
	List(ctx context.Context, filter *models.EquipmentFilter, user *models.User) ([]*models.Equipment, error)
	Create(ctx context.Context, create *models.EquipmentCreate, user *models.User) (*models.Equipment, error)
	Update(ctx context.Context, update *models.EquipmentUpdate, user *models.User) (*models.Equipment, error)
	Delete(ctx context.Context, id string, user *models.User) error
}
