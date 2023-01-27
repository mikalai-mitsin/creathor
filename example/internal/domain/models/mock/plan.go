package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewPlan(t *testing.T) *models.Plan {
	t.Helper()
	return &models.Plan{
		ID:          models.UUID(uuid.NewString()),
		Name:        faker.Lorem().String(),
		Repeat:      uint64(faker.RandomInt(2, 100)),
		EquipmentID: faker.Lorem().String(),
		UpdatedAt:   faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt:   faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewPlanCreate(t *testing.T) *models.PlanCreate {
	t.Helper()
	return &models.PlanCreate{
		Name:        faker.Lorem().String(),
		Repeat:      uint64(faker.RandomInt(2, 100)),
		EquipmentID: faker.Lorem().String(),
	}
}

func NewPlanUpdate(t *testing.T) *models.PlanUpdate {
	t.Helper()
	return &models.PlanUpdate{
		ID:          models.UUID(uuid.NewString()),
		Name:        utils.Pointer(faker.Lorem().String()),
		Repeat:      utils.Pointer(uint64(faker.RandomInt(2, 100))),
		EquipmentID: utils.Pointer(faker.Lorem().String()),
	}
}

func NewPlanFilter(t *testing.T) *models.PlanFilter {
	t.Helper()
	return &models.PlanFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []models.UUID{models.UUID(uuid.NewString()), models.UUID(uuid.NewString()), models.UUID(uuid.NewString())},
		Search:     utils.Pointer(faker.Lorem().String()),
	}
}
