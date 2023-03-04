package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func NewPlan(t *testing.T) *models.Plan {
	t.Helper()
	m := &models.Plan{
		ID:          models.UUID(uuid.NewString()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		Name:        faker.New().Lorem().Sentence(15),
		Repeat:      faker.New().UInt64(),
		EquipmentID: faker.New().Lorem().Sentence(15),
	}
	return m
}
func NewPlanCreate(t *testing.T) *models.PlanCreate {
	t.Helper()
	m := &models.PlanCreate{
		Name:        faker.New().Lorem().Sentence(15),
		Repeat:      faker.New().UInt64(),
		EquipmentID: faker.New().Lorem().Sentence(15),
	}
	return m
}
func NewPlanUpdate(t *testing.T) *models.PlanUpdate {
	t.Helper()
	m := &models.PlanUpdate{
		ID:          models.UUID(uuid.NewString()),
		Name:        utils.Pointer(faker.New().Lorem().Sentence(15)),
		Repeat:      utils.Pointer(faker.New().UInt64()),
		EquipmentID: utils.Pointer(faker.New().Lorem().Sentence(15)),
	}
	return m
}
func NewPlanFilter(t *testing.T) *models.PlanFilter {
	t.Helper()
	m := &models.PlanFilter{
		IDs:        []models.UUID{models.UUID(uuid.NewString()), models.UUID(uuid.NewString())},
		PageNumber: utils.Pointer(faker.New().UInt64()),
		PageSize:   utils.Pointer(faker.New().UInt64()),
		OrderBy:    faker.New().Lorem().Words(27),
		Search:     utils.Pointer(faker.New().Lorem().Sentence(15)),
	}
	return m
}
