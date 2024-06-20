package config

import (
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	DBHost                     string        `env:"DBHost" envDefault:"localhost"`
	DBPort                     int           `env:"DBPort" envDefault:"5432"`
	DBUser                     string        `env:"DBUser" envDefault:"app"`
	DBPassword                 string        `env:"DBPassword" envDefault:"pass"`
	DBName                     string        `env:"DBName" envDefault:"auth"`
	AuthCodeExpires            time.Duration `env:"AuthCodeExpires" envDefault:"120"`
	AuthTokenExpiresMin        time.Duration `env:"AuthTokenExpiresMin" envDefault:"60"`
	AuthRefreshTokenExpiresDay time.Duration `env:"AuthRefreshTokenExpiresDay" envDefault:"30"`
	SessionExpires             int           `env:"SessionExpires" envDefault:"3600"`
}

func GetEnv() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
