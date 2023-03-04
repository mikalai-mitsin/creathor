package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func NewUser(t *testing.T) *models.User {
	t.Helper()
	return &models.User{
		ID:        models.UUID(uuid.NewString()),
		FirstName: faker.New().Person().FirstName(),
		LastName:  faker.New().Person().LastName(),
		Password:  faker.New().Internet().Password(),
		Email:     faker.New().Lorem().Word() + faker.New().Internet().SafeEmail(),
		CreatedAt: faker.New().Time().Time(time.Now()).UTC(),
		UpdatedAt: faker.New().Time().Time(time.Now()).UTC(),
	}
}

func NewUserUpdate(t *testing.T) *models.UserUpdate {
	t.Helper()
	return &models.UserUpdate{
		ID:        models.UUID(uuid.NewString()),
		FirstName: utils.Pointer(faker.New().Person().FirstName()),
		LastName:  utils.Pointer(faker.New().Person().LastName()),
		Password:  utils.Pointer(faker.New().Internet().Password()),
		Email:     utils.Pointer(faker.New().Internet().SafeEmail()),
	}
}

func NewUserFilter(t *testing.T) *models.UserFilter {
	t.Helper()
	return &models.UserFilter{
		PageSize:   utils.Pointer(faker.New().UInt64()),
		PageNumber: utils.Pointer(faker.New().UInt64()),
		Search:     utils.Pointer(faker.New().Lorem().Text(14)),
		OrderBy:    faker.New().Lorem().Words(5),
	}
}

func NewUserCreate(t *testing.T) *models.UserCreate {
	t.Helper()
	return &models.UserCreate{
		Email:    faker.New().Internet().SafeEmail(),
		Password: faker.New().Internet().Password(),
	}
}
