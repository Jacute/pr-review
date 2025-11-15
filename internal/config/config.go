package config

import (
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env string `envconfig:"ENV"`
	*ApplicationConfig
	*DatabaseConfig
}

type ApplicationConfig struct {
	Host         string        `envconfig:"HTTP_HOST"`
	Port         int           `envconfig:"HTTP_PORT"`
	ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `envconfig:"HTTP_IDLE_TIMEOUT"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"POSTGRES_HOST" env-default:"127.0.0.1"`
	Port     int    `envconfig:"POSTGRES_PORT" env-default:"5432"`
	Username string `envconfig:"POSTGRES_USER" env-required:"true" json:"-"`
	Password string `envconfig:"POSTGRES_PASSWORD" env-required:"true" json:"-"`
	Name     string `envconfig:"POSTGRES_DB" env-required:"true" json:"-"`
}

func MustParseConfig() *Config {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "local"
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic("error loading env: " + err.Error())
	}
	cfg.Env = environment

	return &cfg
}
