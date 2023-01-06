package interceptors

import (
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewUserInterceptor, NewUserSessionInterceptor, NewEquipmentInterceptor, NewSessionInterceptor, NewApproachInterceptor),
)
