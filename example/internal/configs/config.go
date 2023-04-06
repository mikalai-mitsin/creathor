package configs

import (
	"github.com/018bf/example/internal/domain/errs"
	"github.com/ilyakaznacheev/cleanenv"
)

type otel struct {
	URL         string `env:"OTEL_URL"         toml:"url"`
	Enabled     bool   `env:"OTEL_ENABLED"     toml:"enabled"`
	Environment string `env:"OTEL_ENVIRONMENT" toml:"environment"`
}

type auth struct {
	PublicKey  string `env:"AUTH_PUBLIC_KEY"  toml:"public_key"`
	PrivateKey string `env:"AUTH_PRIVATE_KEY" toml:"private_key"`
	RefreshTTL int64  `env:"AUTH_REFRESH_TTL" toml:"refresh_ttl" env-default:"172800"`
	AccessTTL  int64  `env:"AUTH_ACCESS_TTL"  toml:"access_ttl"  env-default:"86400"`
}

type database struct {
	URI                string `env:"DATABASE_URI"                  toml:"uri"`
	MaxOpenConnections int    `env:"DATABASE_MAX_OPEN_CONNECTIONS" toml:"max_open_connections" env-default:"50"`
	MaxIDLEConnections int    `env:"DATABASE_MAX_IDLE_CONNECTIONS" toml:"max_idle_connections" env-default:"10"`
}

type Config struct {
	BindAddr string   `env:"BIND_ADDR" toml:"bind_addr" env-default:":8000"`
	LogLevel string   `env:"LOG_LEVEL" toml:"log_level" env-default:"debug"`
	Database database `                toml:"database"`
	Auth     auth     `                toml:"auth"`
	Otel     otel     `                toml:"otel"`
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
