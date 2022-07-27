package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Postgres struct {
		URI string
	}
	HTTP struct {
		Port               int
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxHeaderMegabytes int
	}
}

func New() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("http", &cfg.HTTP); err != nil {
		return nil, err
	}

	if err := envconfig.Process("postgres", &cfg.Postgres); err != nil {
		return nil, err
	}

	return cfg, nil
}
