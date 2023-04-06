package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// EquipmentInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/equipment.go . EquipmentInterceptor
type EquipmentInterceptor interface {
	Get(ctx context.Context, id models.UUID, requestUser *models.User) (*models.Equipment, error)
	List(
		ctx context.Context,
		filter *models.EquipmentFilter,
		requestUser *models.User,
	) ([]*models.Equipment, uint64, error)
	Update(
		ctx context.Context,
		update *models.EquipmentUpdate,
		requestUser *models.User,
	) (*models.Equipment, error)
	Create(
		ctx context.Context,
		create *models.EquipmentCreate,
		requestUser *models.User,
	) (*models.Equipment, error)
	Delete(ctx context.Context, id models.UUID, requestUser *models.User) error
}
