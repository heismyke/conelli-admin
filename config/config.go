package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	HttpPort    int    `envconfig:"HTTP_PORT" default:"8000"`
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
	IsDev       bool   `envconfig:"IS_DEV" default:"true"`
	CorsOrigin  string `envconfig:"CORS_ORIGIN" default:"http://localhost:5173"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
