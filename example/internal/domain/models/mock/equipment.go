package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func NewEquipment(t *testing.T) *models.Equipment {
	t.Helper()
	m := &models.Equipment{
		ID:        models.UUID(uuid.NewString()),
		UpdatedAt: faker.New().Time().Time(time.Now()),
		CreatedAt: faker.New().Time().Time(time.Now()),
		Name:      faker.New().Lorem().Sentence(15),
		Repeat:    faker.New().Int(),
		Weight:    faker.New().Int(),
	}
	return m
}
func NewEquipmentCreate(t *testing.T) *models.EquipmentCreate {
	t.Helper()
	m := &models.EquipmentCreate{
		Name:   faker.New().Lorem().Sentence(15),
		Repeat: faker.New().Int(),
		Weight: faker.New().Int(),
	}
	return m
}
func NewEquipmentUpdate(t *testing.T) *models.EquipmentUpdate {
	t.Helper()
	m := &models.EquipmentUpdate{
		ID:     models.UUID(uuid.NewString()),
		Name:   utils.Pointer(faker.New().Lorem().Sentence(15)),
		Repeat: utils.Pointer(faker.New().Int()),
		Weight: utils.Pointer(faker.New().Int()),
	}
	return m
}
func NewEquipmentFilter(t *testing.T) *models.EquipmentFilter {
	t.Helper()
	m := &models.EquipmentFilter{
		IDs:        []models.UUID{models.UUID(uuid.NewString()), models.UUID(uuid.NewString())},
		PageNumber: utils.Pointer(faker.New().UInt64()),
		PageSize:   utils.Pointer(faker.New().UInt64()),
		OrderBy:    faker.New().Lorem().Words(27),
		Search:     utils.Pointer(faker.New().Lorem().Sentence(15)),
	}
	return m
}
