package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type App struct {
	Port      string `env:"PORT" envDefault:"8080"`
	WebOrigin string `env:"WEB_ORIGIN"`
}

type Postgres struct {
	URL string `env:"PG_URL,required"`
}

type Redis struct {
	Addr     string `env:"REDIS_ADDR,required"`
	Username string `env:"REDIS_USERNAME"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type Storage struct {
	BaseURL  string `env:"STORAGE_BASE_URL" envDefault:"http://localhost:8080"`
	BasePath string `env:"STORAGE_BASE_PATH" envDefault:"./uploads/"`
}

type Config struct {
	App      App
	Postgres Postgres
	Redis    Redis
	Storage  Storage
}

func New() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse environment variables: %v", err)
	}

	return cfg, nil
}
