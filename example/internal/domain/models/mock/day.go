package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func NewDay(t *testing.T) *models.Day {
	t.Helper()
	m := &models.Day{
		ID:          models.UUID(uuid.NewString()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		Name:        faker.New().Lorem().Sentence(15),
		Repeat:      faker.New().Int(),
		EquipmentID: faker.New().Lorem().Sentence(15),
	}
	return m
}
func NewDayCreate(t *testing.T) *models.DayCreate {
	t.Helper()
	m := &models.DayCreate{
		Name:        faker.New().Lorem().Sentence(15),
		Repeat:      faker.New().Int(),
		EquipmentID: faker.New().Lorem().Sentence(15),
	}
	return m
}
func NewDayUpdate(t *testing.T) *models.DayUpdate {
	t.Helper()
	m := &models.DayUpdate{
		ID:          models.UUID(uuid.NewString()),
		Name:        utils.Pointer(faker.New().Lorem().Sentence(15)),
		Repeat:      utils.Pointer(faker.New().Int()),
		EquipmentID: utils.Pointer(faker.New().Lorem().Sentence(15)),
	}
	return m
}
func NewDayFilter(t *testing.T) *models.DayFilter {
	t.Helper()
	m := &models.DayFilter{
		IDs:        []models.UUID{models.UUID(uuid.NewString()), models.UUID(uuid.NewString())},
		PageNumber: utils.Pointer(faker.New().UInt64()),
		PageSize:   utils.Pointer(faker.New().UInt64()),
		OrderBy:    faker.New().Lorem().Words(27),
		Search:     utils.Pointer(faker.New().Lorem().Sentence(15)),
	}
	return m
}
