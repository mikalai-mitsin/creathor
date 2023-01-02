package mock_models

import (
	"github.com/018bf/creathor/internal/domain/models"
	"github.com/018bf/creathor/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
	"testing"
)

func NewUserApproach(t *testing.T) *models.UserApproach {
	t.Helper()
	return &models.UserApproach{
		ID: uuid.New().String(),
	}
}

func NewUserApproachCreate(t *testing.T) *models.UserApproachCreate {
	t.Helper()
	return &models.UserApproachCreate{}
}

func NewUserApproachUpdate(t *testing.T) *models.UserApproachUpdate {
	t.Helper()
	return &models.UserApproachUpdate{
		ID: uuid.New().String(),
	}
}

func NewUserApproachFilter(t *testing.T) *models.UserApproachFilter {
	t.Helper()
	return &models.UserApproachFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
}
