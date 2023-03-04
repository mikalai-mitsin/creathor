package mock_models

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

func NewArch(t *testing.T) *models.Arch {
	t.Helper()
	m := &models.Arch{
		ID:          models.UUID(uuid.NewString()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		Name:        faker.New().Lorem().Sentence(15),
		Title:       faker.New().Lorem().Sentence(15),
		Description: faker.New().Lorem().Sentence(15),
		Tags:        faker.New().Lorem().Words(27),
		Versions:    []uint64{faker.New().UInt64(), faker.New().UInt64()},
		Release:     faker.New().Time().Time(time.Now()),
		Tested:      faker.New().Time().Time(time.Now()),
	}
	return m
}
func NewArchCreate(t *testing.T) *models.ArchCreate {
	t.Helper()
	m := &models.ArchCreate{
		Name:        faker.New().Lorem().Sentence(15),
		Title:       faker.New().Lorem().Sentence(15),
		Description: faker.New().Lorem().Sentence(15),
		Tags:        faker.New().Lorem().Words(27),
		Versions:    []uint64{faker.New().UInt64(), faker.New().UInt64()},
		Release:     faker.New().Time().Time(time.Now()),
		Tested:      faker.New().Time().Time(time.Now()),
	}
	return m
}
func NewArchUpdate(t *testing.T) *models.ArchUpdate {
	t.Helper()
	m := &models.ArchUpdate{
		ID:          models.UUID(uuid.NewString()),
		Name:        utils.Pointer(faker.New().Lorem().Sentence(15)),
		Title:       utils.Pointer(faker.New().Lorem().Sentence(15)),
		Description: utils.Pointer(faker.New().Lorem().Sentence(15)),
		Tags:        utils.Pointer(faker.New().Lorem().Words(27)),
		Versions:    utils.Pointer([]uint64{faker.New().UInt64(), faker.New().UInt64()}),
		Release:     utils.Pointer(faker.New().Time().Time(time.Now())),
		Tested:      utils.Pointer(faker.New().Time().Time(time.Now())),
	}
	return m
}
func NewArchFilter(t *testing.T) *models.ArchFilter {
	t.Helper()
	m := &models.ArchFilter{
		IDs:        []models.UUID{models.UUID(uuid.NewString()), models.UUID(uuid.NewString())},
		PageNumber: utils.Pointer(faker.New().UInt64()),
		PageSize:   utils.Pointer(faker.New().UInt64()),
		OrderBy:    faker.New().Lorem().Words(27),
		Search:     utils.Pointer(faker.New().Lorem().Sentence(15)),
	}
	return m
}
