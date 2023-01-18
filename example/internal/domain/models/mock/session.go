package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewSession(t *testing.T) *models.Session {
	t.Helper()
	return &models.Session{
		ID:          uuid.New().String(),
		Description: faker.Lorem().String(),
		Title:       faker.Lorem().String(),
		UpdatedAt:   faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt:   faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewSessionCreate(t *testing.T) *models.SessionCreate {
	t.Helper()
	return &models.SessionCreate{
		Description: faker.Lorem().String(),
		Title:       faker.Lorem().String(),
	}
}

func NewSessionUpdate(t *testing.T) *models.SessionUpdate {
	t.Helper()
	return &models.SessionUpdate{
		ID:          uuid.New().String(),
		Description: utils.Pointer(faker.Lorem().String()),
		Title:       utils.Pointer(faker.Lorem().String()),
	}
}

func NewSessionFilter(t *testing.T) *models.SessionFilter {
	t.Helper()
	return &models.SessionFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
}
