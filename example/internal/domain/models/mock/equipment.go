package mock_models // nolint:stylecheck

import (
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
	"testing"
	"time"
)

func NewEquipment(t *testing.T) *models.Equipment {
	t.Helper()
	return &models.Equipment{
		ID:        uuid.New().String(),
		UpdatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewEquipmentCreate(t *testing.T) *models.EquipmentCreate {
	t.Helper()
	return &models.EquipmentCreate{}
}

func NewEquipmentUpdate(t *testing.T) *models.EquipmentUpdate {
	t.Helper()
	return &models.EquipmentUpdate{
		ID: uuid.New().String(),
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
