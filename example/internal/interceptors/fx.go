package interceptors

import (
	"go.uber.org/fx"
)

var FXModule = fx.Options(
	fx.Provide(NewAuthInterceptor, NewUserInterceptor, NewUserSessionInterceptor, NewEquipmentInterceptor, NewSessionInterceptor, NewApproachInterceptor, NewMarkInterceptor),
)
