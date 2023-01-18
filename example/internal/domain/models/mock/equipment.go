package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewEquipment(t *testing.T) *models.Equipment {
	t.Helper()
	return &models.Equipment{
		ID:        uuid.New().String(),
		Name:      faker.Lorem().String(),
		Repeat:    faker.RandomInt(2, 100),
		Weight:    faker.RandomInt(2, 100),
		UpdatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewEquipmentCreate(t *testing.T) *models.EquipmentCreate {
	t.Helper()
	return &models.EquipmentCreate{
		Name:   faker.Lorem().String(),
		Repeat: faker.RandomInt(2, 100),
		Weight: faker.RandomInt(2, 100),
	}
}

func NewEquipmentUpdate(t *testing.T) *models.EquipmentUpdate {
	t.Helper()
	return &models.EquipmentUpdate{
		ID:     uuid.New().String(),
		Name:   utils.Pointer(faker.Lorem().String()),
		Repeat: utils.Pointer(faker.RandomInt(2, 100)),
		Weight: utils.Pointer(faker.RandomInt(2, 100)),
	}
}

func NewEquipmentFilter(t *testing.T) *models.EquipmentFilter {
	t.Helper()
	return &models.EquipmentFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
}
