package repositories

import (
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewJWTAuthRepository, NewPostgresUserRepository, NewPermissionRepository, NewUserSessionRepository, NewEquipmentRepository, NewSessionRepository, NewApproachRepository, NewMarkRepository),
)
