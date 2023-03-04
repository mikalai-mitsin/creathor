package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

type ArchInterceptor struct {
	archUseCase usecases.ArchUseCase
	logger      log.Logger
	authUseCase usecases.AuthUseCase
}

func NewArchInterceptor(
	archUseCase usecases.ArchUseCase,
	logger log.Logger,
	authUseCase usecases.AuthUseCase,
) interceptors.ArchInterceptor {
	return &ArchInterceptor{archUseCase: archUseCase, logger: logger, authUseCase: authUseCase}
}

func (i *ArchInterceptor) Create(
	ctx context.Context,
	create *models.ArchCreate,
	requestUser *models.User,
) (*models.Arch, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDArchCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDArchCreate, create); err != nil {
		return nil, err
	}
	arch, err := i.archUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return arch, nil
}

func (i *ArchInterceptor) Get(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) (*models.Arch, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDArchDetail); err != nil {
		return nil, err
	}
	arch, err := i.archUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDArchDetail, arch); err != nil {
		return nil, err
	}
	return arch, nil
}

func (i *ArchInterceptor) List(
	ctx context.Context,
	filter *models.ArchFilter,
	requestUser *models.User,
) ([]*models.Arch, uint64, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDArchList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDArchList, filter); err != nil {
		return nil, 0, err
	}
	listArches, count, err := i.archUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return listArches, count, nil
}

func (i *ArchInterceptor) Update(
	ctx context.Context,
	update *models.ArchUpdate,
	requestUser *models.User,
) (*models.Arch, error) {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDArchUpdate); err != nil {
		return nil, err
	}
	arch, err := i.archUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDArchUpdate, arch); err != nil {
		return nil, err
	}
	updated, err := i.archUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (i *ArchInterceptor) Delete(
	ctx context.Context,
	id models.UUID,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(ctx, requestUser, models.PermissionIDArchDelete); err != nil {
		return err
	}
	arch, err := i.archUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, models.PermissionIDArchDelete, arch); err != nil {
		return err
	}
	if err := i.archUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
