package repositories

import (
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewUserRepository, NewUserSessionRepository, NewEquipmentRepository, NewSessionRepository, NewApproachRepository),
)
