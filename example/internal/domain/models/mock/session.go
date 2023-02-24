package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func NewSession(t *testing.T) *models.Session {
	t.Helper()
	m := &models.Session{
		ID:          models.UUID(uuid.NewString()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		Title:       faker.New().Lorem().Text(256),
		Description: faker.New().Lorem().Text(256),
		Weight:      faker.New().UInt64(),
		Versions:    []uint64{faker.New().UInt64(), faker.New().UInt64()},
		Release:     faker.New().Time().Time(time.Now()),
		Tested:      faker.New().Time().Time(time.Now()),
	}
	return m
}
func NewSessionCreate(t *testing.T) *models.SessionCreate {
	t.Helper()
	m := &models.SessionCreate{
		Title:       faker.New().Lorem().Text(256),
		Description: faker.New().Lorem().Text(256),
		Weight:      faker.New().UInt64(),
		Versions:    []uint64{faker.New().UInt64(), faker.New().UInt64()},
		Release:     faker.New().Time().Time(time.Now()),
		Tested:      faker.New().Time().Time(time.Now()),
	}
	return m
}
func NewSessionUpdate(t *testing.T) *models.SessionUpdate {
	t.Helper()
	m := &models.SessionUpdate{
		ID:          models.UUID(uuid.NewString()),
		Title:       utils.Pointer(faker.New().Lorem().Text(256)),
		Description: utils.Pointer(faker.New().Lorem().Text(256)),
		Weight:      utils.Pointer(faker.New().UInt64()),
		Versions:    utils.Pointer([]uint64{faker.New().UInt64(), faker.New().UInt64()}),
		Release:     utils.Pointer(faker.New().Time().Time(time.Now())),
		Tested:      utils.Pointer(faker.New().Time().Time(time.Now())),
	}
	return m
}
func NewSessionFilter(t *testing.T) *models.SessionFilter {
	t.Helper()
	m := &models.SessionFilter{
		IDs:        []models.UUID{models.UUID(uuid.NewString()), models.UUID(uuid.NewString())},
		PageNumber: utils.Pointer(faker.New().UInt64()),
		PageSize:   utils.Pointer(faker.New().UInt64()),
		OrderBy:    faker.New().Lorem().Words(27),
		Search:     utils.Pointer(faker.New().Lorem().Text(256)),
	}
	return m
}
