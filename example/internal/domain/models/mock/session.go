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
		Title:       faker.New().Lorem().Sentence(15),
		Description: faker.New().Lorem().Sentence(15),
	}
	return m
}
func NewSessionCreate(t *testing.T) *models.SessionCreate {
	t.Helper()
	m := &models.SessionCreate{
		Title:       faker.New().Lorem().Sentence(15),
		Description: faker.New().Lorem().Sentence(15),
	}
	return m
}
func NewSessionUpdate(t *testing.T) *models.SessionUpdate {
	t.Helper()
	m := &models.SessionUpdate{
		ID:          models.UUID(uuid.NewString()),
		Title:       utils.Pointer(faker.New().Lorem().Sentence(15)),
		Description: utils.Pointer(faker.New().Lorem().Sentence(15)),
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
		Search:     utils.Pointer(faker.New().Lorem().Sentence(15)),
	}
	return m
}
