package mock_models // nolint:stylecheck

import (
	"testing"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewApproach(t *testing.T) *models.Approach {
	t.Helper()
	return &models.Approach{
		ID: uuid.New().String(),
	}
}

func NewApproachCreate(t *testing.T) *models.ApproachCreate {
	t.Helper()
	return &models.ApproachCreate{}
}

func NewApproachUpdate(t *testing.T) *models.ApproachUpdate {
	t.Helper()
	return &models.ApproachUpdate{
		ID: uuid.New().String(),
	}
}

func NewApproachFilter(t *testing.T) *models.ApproachFilter {
	t.Helper()
	return &models.ApproachFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
}
