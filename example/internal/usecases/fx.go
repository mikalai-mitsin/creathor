package usecases

import (
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewUserUseCase, NewUserSessionUseCase, NewEquipmentUseCase, NewSessionUseCase, NewApproachUseCase),
)
