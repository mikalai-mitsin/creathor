package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewUser(t *testing.T) *models.User {
	t.Helper()
	return &models.User{
		ID:        uuid.New().String(),
		UpdatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewUserCreate(t *testing.T) *models.UserCreate {
	t.Helper()
	return &models.UserCreate{}
}

func NewUserUpdate(t *testing.T) *models.UserUpdate {
	t.Helper()
	return &models.UserUpdate{
		ID: uuid.New().String(),
	}
}

func NewUserFilter(t *testing.T) *models.UserFilter {
	t.Helper()
	return &models.UserFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
}
