package interceptors

import (
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewAuthInterceptor, NewUserInterceptor, NewSessionInterceptor, NewEquipmentInterceptor),
)
