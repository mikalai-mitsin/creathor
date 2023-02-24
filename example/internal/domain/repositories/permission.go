package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

//nolint: lll
//go:generate mockgen -build_flags=-mod=mod -build_flags=-mod=mod -destination mock/permission_mock.go github.com/018bf/example/internal/domain/repositories PermissionRepository

type PermissionRepository interface {
	HasPermission(
		ctx context.Context,
		permission models.PermissionID,
		requestUser *models.User,
	) error
	HasObjectPermission(
		ctx context.Context,
		permission models.PermissionID,
		user *models.User,
		obj any,
	) error
}
