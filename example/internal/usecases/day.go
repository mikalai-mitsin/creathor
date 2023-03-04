package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type DayUseCase struct {
	dayRepository repositories.DayRepository
	clock         clock.Clock
	logger        log.Logger
}

func NewDayUseCase(
	dayRepository repositories.DayRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.DayUseCase {
	return &DayUseCase{dayRepository: dayRepository, clock: clock, logger: logger}
}
func (u *DayUseCase) Create(ctx context.Context, create *models.DayCreate) (*models.Day, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	day := &models.Day{
		ID:          "",
		UpdatedAt:   now,
		CreatedAt:   now,
		Name:        create.Name,
		Repeat:      create.Repeat,
		EquipmentID: create.EquipmentID,
	}
	if err := u.dayRepository.Create(ctx, day); err != nil {
		return nil, err
	}
	return day, nil
}
func (u *DayUseCase) Get(ctx context.Context, id models.UUID) (*models.Day, error) {
	day, err := u.dayRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return day, nil
}

func (u *DayUseCase) List(
	ctx context.Context,
	filter *models.DayFilter,
) ([]*models.Day, uint64, error) {
	day, err := u.dayRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.dayRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return day, count, nil
}
func (u *DayUseCase) Update(ctx context.Context, update *models.DayUpdate) (*models.Day, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	day, err := u.dayRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	{
		if update.Name != nil {
			day.Name = *update.Name
		}
		if update.Repeat != nil {
			day.Repeat = *update.Repeat
		}
		if update.EquipmentID != nil {
			day.EquipmentID = *update.EquipmentID
		}
	}
	day.UpdatedAt = u.clock.Now().UTC()
	if err := u.dayRepository.Update(ctx, day); err != nil {
		return nil, err
	}
	return day, nil
}
func (u *DayUseCase) Delete(ctx context.Context, id models.UUID) error {
	if err := u.dayRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
