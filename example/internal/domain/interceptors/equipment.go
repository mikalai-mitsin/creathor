package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -destination mock/equipment_mock.go github.com/018bf/example/internal/domain/interceptors EquipmentInterceptor

type EquipmentInterceptor interface {
	Get(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) (*models.Equipment, error)
	List(
		ctx context.Context,
		filter *models.EquipmentFilter,
		requestUser *models.User,
	) ([]*models.Equipment, uint64, error)
	Create(
		ctx context.Context,
		create *models.EquipmentCreate,
		requestUser *models.User,
	) (*models.Equipment, error)
	Update(
		ctx context.Context,
		update *models.EquipmentUpdate,
		requestUser *models.User,
	) (*models.Equipment, error)
	Delete(
		ctx context.Context,
		id string,
		requestUser *models.User,
	) error
}
