package cmd

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Database     DatabaseConfig
	Cache        CacheConfig
	Microservice MicroserviceConfig
}

type DatabaseConfig struct {
	Host     string `env:"PG_HOST,required"`
	Port     string `env:"PG_PORT,required"`
	User     string `env:"PG_USER,required"`
	Password string `env:"PG_PASSWORD,required"`
	Name     string `env:"PG_DATABASE_NAME,required"`
}

type CacheConfig struct {
	Host     string `env:"CACHE_HOST,required"`
	Port     string `env:"CACHE_PORT,required"`
}

type MicroserviceConfig struct {
	Port string `env:"MICROSERVICE_PORT,required"`
}

func Load() (*Config, error) {
	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
