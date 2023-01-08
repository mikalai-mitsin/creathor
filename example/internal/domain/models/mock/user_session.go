package mock_models // nolint:stylecheck

import (
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
	"testing"
	"time"
)

func NewUserSession(t *testing.T) *models.UserSession {
	t.Helper()
	return &models.UserSession{
		ID:        uuid.New().String(),
		UpdatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewUserSessionCreate(t *testing.T) *models.UserSessionCreate {
	t.Helper()
	return &models.UserSessionCreate{}
}

func NewUserSessionUpdate(t *testing.T) *models.UserSessionUpdate {
	t.Helper()
	return &models.UserSessionUpdate{
		ID: uuid.New().String(),
	}
}

func NewUserSessionFilter(t *testing.T) *models.UserSessionFilter {
	t.Helper()
	return &models.UserSessionFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
}
