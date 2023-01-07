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
		ID:        uuid.NewString(),
		FirstName: faker.Name().FirstName(),
		LastName:  faker.Name().LastName(),
		Password:  faker.Number().Hexadecimal(10),
		Email:     faker.Lorem().Word() + faker.Internet().SafeEmail(),
		CreatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
		UpdatedAt: faker.Time().Backward(time.Hour).UTC(),
	}
}

func NewUserUpdate(t *testing.T) *models.UserUpdate {
	t.Helper()
	return &models.UserUpdate{
		ID:        uuid.NewString(),
		FirstName: utils.Pointer(faker.Name().FirstName()),
		LastName:  utils.Pointer(faker.Name().LastName()),
		Password:  utils.Pointer(faker.Number().Hexadecimal(10)),
		Email:     utils.Pointer(faker.Internet().SafeEmail()),
	}
}

func NewUserFilter(t *testing.T) *models.UserFilter {
	t.Helper()
	return &models.UserFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		Search:     utils.Pointer(faker.Lorem().String()),
		OrderBy:    faker.Lorem().Words(5),
	}
}

func NewUserCreate(t *testing.T) *models.UserCreate {
	t.Helper()
	return &models.UserCreate{
		Email:    faker.Internet().SafeEmail(),
		Password: faker.Internet().Password(6, 12),
	}
}
