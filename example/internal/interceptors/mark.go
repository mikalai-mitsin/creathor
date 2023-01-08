package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type MarkInterceptor struct {
	markUseCase usecases.MarkUseCase
	authUseCase usecases.AuthUseCase
	logger      log.Logger
}

func NewMarkInterceptor(
	markUseCase usecases.MarkUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.MarkInterceptor {
	return &MarkInterceptor{
		markUseCase: markUseCase,
		authUseCase: authUseCase,
		logger:      logger,
	}
}

func (i *MarkInterceptor) Get(
	ctx context.Context,
	id string,
	requestUser *models.User,
) (*models.Mark, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkDetail,
	); err != nil {
		return nil, err
	}
	mark, err := i.markUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkDetail,
		mark,
	)
	if err != nil {
		return nil, err
	}
	return mark, nil
}

func (i *MarkInterceptor) List(
	ctx context.Context,
	filter *models.MarkFilter,
	requestUser *models.User,
) ([]*models.Mark, uint64, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkList,
	); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkList,
		filter,
	); err != nil {
		return nil, 0, err
	}
	marks, count, err := i.markUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return marks, count, nil
}

func (i *MarkInterceptor) Create(
	ctx context.Context,
	create *models.MarkCreate,
	requestUser *models.User,
) (*models.Mark, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkCreate,
	); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkCreate,
		create,
	); err != nil {
		return nil, err
	}
	mark, err := i.markUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return mark, nil
}

func (i *MarkInterceptor) Update(
	ctx context.Context,
	update *models.MarkUpdate,
	requestUser *models.User,
) (*models.Mark, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkUpdate,
	); err != nil {
		return nil, err
	}
	mark, err := i.markUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkUpdate,
		mark,
	); err != nil {
		return nil, err
	}
	updatedMark, err := i.markUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedMark, nil
}

func (i *MarkInterceptor) Delete(
	ctx context.Context,
	id string,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkDelete,
	); err != nil {
		return err
	}
	mark, err := i.markUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDMarkDelete,
		mark,
	)
	if err != nil {
		return err
	}
	if err := i.markUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
