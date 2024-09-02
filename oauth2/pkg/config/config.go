package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	DBHost                     string `env:"DBHost" envDefault:"localhost"`
	DBPort                     uint16 `env:"DBPort" envDefault:"5432"`
	DBUser                     string `env:"DBUser" envDefault:"app"`
	DBPassword                 string `env:"DBPassword" envDefault:"pass"`
	DBName                     string `env:"DBName" envDefault:"auth"`
	AuthCodeExpires            int    `env:"AuthCodeExpires" envDefault:"120"`           // 秒を単位として指定
	AuthTokenExpiresMin        int    `env:"AuthTokenExpiresMin" envDefault:"60"`        // 分を単位として指定
	AuthRefreshTokenExpiresDay int    `env:"AuthRefreshTokenExpiresDay" envDefault:"30"` // 時間を単位として指定
	SessionExpires             int    `env:"SessionExpires" envDefault:"3600"`
}

func GetEnv() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
