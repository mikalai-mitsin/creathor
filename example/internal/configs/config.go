package configs

import (
	"github.com/018bf/example/internal/domain/errs"
	"github.com/ilyakaznacheev/cleanenv"
)

type auth struct {
	PublicKey  string `env:"AUTH_PUBLIC_KEY" toml:"public_key"`
	PrivateKey string `env:"AUTH_PRIVATE_KEY" toml:"private_key"`
	RefreshTTL int64  `env:"AUTH_REFRESH_TTL"  env-default:"172800"  toml:"refresh_ttl"`
	AccessTTL  int64  `env:"AUTH_ACCESS_TTL"  env-default:"86400"  toml:"access_ttl"`
}

type database struct {
	URI                string `env:"DATABASE_URI" toml:"uri"`
	MaxOpenConnections int    `env:"DATABASE_MAX_OPEN_CONNECTIONS" env-default:"50"  toml:"max_open_connections"`
	MaxIDLEConnections int    `env:"DATABASE_MAX_IDLE_CONNECTIONS" env-default:"10"  toml:"max_idle_connections"`
}

type Config struct {
	BindAddr string   `env:"BIND_ADDR" toml:"bind_addr" env-default:":8000"`
	LogLevel string   `env:"LOG_LEVEL" toml:"log_level" env-default:"debug"`
	Database database `toml:"database"`
	Auth     auth     `toml:"auth"`
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
