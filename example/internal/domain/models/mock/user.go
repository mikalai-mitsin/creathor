package mock_models

import (
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
	"testing"
)

func NewUser(t *testing.T) *models.User {
	t.Helper()
	return &models.User{
		ID: uuid.New().String(),
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
