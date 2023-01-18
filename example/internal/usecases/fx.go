package usecases

import (
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewAuthUseCase, NewUserUseCase, NewSessionUseCase, NewEquipmentUseCase),
)
