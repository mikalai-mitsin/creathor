package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewDay(t *testing.T) *models.Day {
	t.Helper()
	return &models.Day{
		ID:          models.UUID(uuid.NewString()),
		Name:        faker.Lorem().String(),
		Repeat:      faker.RandomInt(2, 100),
		EquipmentID: faker.Lorem().String(),
		UpdatedAt:   faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt:   faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewDayCreate(t *testing.T) *models.DayCreate {
	t.Helper()
	return &models.DayCreate{
		Name:        faker.Lorem().String(),
		Repeat:      faker.RandomInt(2, 100),
		EquipmentID: faker.Lorem().String(),
	}
}

func NewDayUpdate(t *testing.T) *models.DayUpdate {
	t.Helper()
	return &models.DayUpdate{
		ID:          models.UUID(uuid.NewString()),
		Name:        utils.Pointer(faker.Lorem().String()),
		Repeat:      utils.Pointer(faker.RandomInt(2, 100)),
		EquipmentID: utils.Pointer(faker.Lorem().String()),
	}
}

func NewDayFilter(t *testing.T) *models.DayFilter {
	t.Helper()
	return &models.DayFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs: []models.UUID{
			models.UUID(uuid.NewString()),
			models.UUID(uuid.NewString()),
			models.UUID(uuid.NewString()),
		},
		Search: utils.Pointer(faker.Lorem().String()),
	}
}
