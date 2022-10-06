package configs

import (
	"github.com/018bf/example/internal/domain/errs"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BindAddr string `env:"BIND_ADDR" toml:"bind_addr" env-default:":8000"`
	LogLevel string `env:"LOG_LEVEL" toml:"log_level" env-default:"debug"`
}

func ParseConfig(configPath string) (*Config, error) {
	config := &Config{}
	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, config); err != nil {
			return nil, errs.NewUnexpectedBehaviorError(err.Error())
		}
	} else {
		if err := cleanenv.ReadEnv(config); err != nil {
			return nil, errs.NewUnexpectedBehaviorError(err.Error())
		}
	}
	return config, nil
}
