package usecases

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/clock"
	"github.com/018bf/example/pkg/log"
)

type ArchUseCase struct {
	archRepository repositories.ArchRepository
	clock          clock.Clock
	logger         log.Logger
}

func NewArchUseCase(
	archRepository repositories.ArchRepository,
	clock clock.Clock,
	logger log.Logger,
) usecases.ArchUseCase {
	return &ArchUseCase{
		archRepository: archRepository,
		clock:          clock,
		logger:         logger,
	}
}

func (u *ArchUseCase) Get(
	ctx context.Context,
	id models.UUID,
) (*models.Arch, error) {
	arch, err := u.archRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return arch, nil
}

func (u *ArchUseCase) List(
	ctx context.Context,
	filter *models.ArchFilter,
) ([]*models.Arch, uint64, error) {
	listArches, err := u.archRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.archRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return listArches, count, nil
}

func (u *ArchUseCase) Create(
	ctx context.Context,
	create *models.ArchCreate,
) (*models.Arch, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	arch := &models.Arch{
		ID:          "",
		Name:        create.Name,
		Tags:        create.Tags,
		Versions:    create.Versions,
		OldVersions: create.OldVersions,
		Release:     create.Release,
		Tested:      create.Tested,
		UpdatedAt:   now,
		CreatedAt:   now,
	}
	if err := u.archRepository.Create(ctx, arch); err != nil {
		return nil, err
	}
	return arch, nil
}

func (u *ArchUseCase) Update(
	ctx context.Context,
	update *models.ArchUpdate,
) (*models.Arch, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	arch, err := u.archRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if update.Name != nil {
		arch.Name = *update.Name
	}
	if update.Tags != nil {
		arch.Tags = *update.Tags
	}
	if update.Versions != nil {
		arch.Versions = *update.Versions
	}
	if update.OldVersions != nil {
		arch.OldVersions = *update.OldVersions
	}
	if update.Release != nil {
		arch.Release = *update.Release
	}
	if update.Tested != nil {
		arch.Tested = *update.Tested
	}
	arch.UpdatedAt = u.clock.Now()
	if err := u.archRepository.Update(ctx, arch); err != nil {
		return nil, err
	}
	return arch, nil
}

func (u *ArchUseCase) Delete(ctx context.Context, id models.UUID) error {
	if err := u.archRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
