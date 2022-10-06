package configs

import (
	"go.uber.org/fx"
)

var FXModule = fx.Provide(ParseConfig)
