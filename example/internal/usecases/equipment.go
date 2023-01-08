package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type EquipmentUseCase struct {
	equipmentRepository repositories.EquipmentRepository
	clock               clock.Clock
	logger              log.Logger
}

func NewEquipmentUseCase(
	equipmentRepository repositories.EquipmentRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.EquipmentUseCase {
	return &EquipmentUseCase{
		equipmentRepository: equipmentRepository,
		clock:               clock,
		logger:              logger,
	}
}

func (u *EquipmentUseCase) Get(
	ctx context.Context,
	id string,
) (*models.Equipment, error) {
	equipment, err := u.equipmentRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (u *EquipmentUseCase) List(
	ctx context.Context,
	filter *models.EquipmentFilter,
) ([]*models.Equipment, uint64, error) {
	equipment, err := u.equipmentRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.equipmentRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return equipment, count, nil
}

func (u *EquipmentUseCase) Create(
	ctx context.Context,
	create *models.EquipmentCreate,
) (*models.Equipment, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	equipment := &models.Equipment{
		ID:        "",
		UpdatedAt: now,
		CreatedAt: now,
	}
	if err := u.equipmentRepository.Create(ctx, equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (u *EquipmentUseCase) Update(
	ctx context.Context,
	update *models.EquipmentUpdate,
) (*models.Equipment, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	equipment, err := u.equipmentRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	equipment.UpdatedAt = u.clock.Now()
	if err := u.equipmentRepository.Update(ctx, equipment); err != nil {
		return nil, err
	}
	return equipment, nil
}

func (u *EquipmentUseCase) Delete(ctx context.Context, id string) error {
	if err := u.equipmentRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
