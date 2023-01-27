package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewArch(t *testing.T) *models.Arch {
	t.Helper()
	return &models.Arch{
		ID:        models.UUID(uuid.NewString()),
		Name:      faker.Lorem().String(),
		Release:   faker.Time().Backward(40 * time.Hour).UTC(),
		Tested:    faker.Time().Backward(40 * time.Hour).UTC(),
		UpdatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewArchCreate(t *testing.T) *models.ArchCreate {
	t.Helper()
	return &models.ArchCreate{
		Name:    faker.Lorem().String(),
		Release: faker.Time().Backward(40 * time.Hour).UTC(),
		Tested:  faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewArchUpdate(t *testing.T) *models.ArchUpdate {
	t.Helper()
	return &models.ArchUpdate{
		ID:      models.UUID(uuid.NewString()),
		Name:    utils.Pointer(faker.Lorem().String()),
		Release: utils.Pointer(faker.Time().Backward(40 * time.Hour).UTC()),
		Tested:  utils.Pointer(faker.Time().Backward(40 * time.Hour).UTC()),
	}
}

func NewArchFilter(t *testing.T) *models.ArchFilter {
	t.Helper()
	return &models.ArchFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []models.UUID{models.UUID(uuid.NewString()), models.UUID(uuid.NewString()), models.UUID(uuid.NewString())},
		Search:     utils.Pointer(faker.Lorem().String()),
	}
}
