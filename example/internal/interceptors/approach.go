package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type ApproachInterceptor struct {
	approachUseCase usecases.ApproachUseCase
	authUseCase     usecases.AuthUseCase
	logger          log.Logger
}

func NewApproachInterceptor(
	approachUseCase usecases.ApproachUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.ApproachInterceptor {
	return &ApproachInterceptor{
		approachUseCase: approachUseCase,
		authUseCase:     authUseCase,
		logger:          logger,
	}
}

func (i *ApproachInterceptor) Get(
	ctx context.Context,
	id string,
	requestUser *models.User,
) (*models.Approach, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachDetail,
	); err != nil {
		return nil, err
	}
	approach, err := i.approachUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachDetail,
		approach,
	)
	if err != nil {
		return nil, err
	}
	return approach, nil
}

func (i *ApproachInterceptor) List(
	ctx context.Context,
	filter *models.ApproachFilter,
	requestUser *models.User,
) ([]*models.Approach, uint64, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachList,
	); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachList,
		filter,
	); err != nil {
		return nil, 0, err
	}
	approaches, count, err := i.approachUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return approaches, count, nil
}

func (i *ApproachInterceptor) Create(
	ctx context.Context,
	create *models.ApproachCreate,
	requestUser *models.User,
) (*models.Approach, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachCreate,
	); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachCreate,
		create,
	); err != nil {
		return nil, err
	}
	approach, err := i.approachUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return approach, nil
}

func (i *ApproachInterceptor) Update(
	ctx context.Context,
	update *models.ApproachUpdate,
	requestUser *models.User,
) (*models.Approach, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachUpdate,
	); err != nil {
		return nil, err
	}
	approach, err := i.approachUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachUpdate,
		approach,
	); err != nil {
		return nil, err
	}
	updatedApproach, err := i.approachUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedApproach, nil
}

func (i *ApproachInterceptor) Delete(
	ctx context.Context,
	id string,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachDelete,
	); err != nil {
		return err
	}
	approach, err := i.approachUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDApproachDelete,
		approach,
	)
	if err != nil {
		return err
	}
	if err := i.approachUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
