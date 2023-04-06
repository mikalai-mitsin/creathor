package repositories

import (
	"context"

	"github.com/018bf/example/internal/domain/models"
)

// PermissionRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/permission.go . PermissionRepository
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
