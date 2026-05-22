package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port   string `env:"PORT" env-default:"8080"`
	AppEnv string `env:"APP_ENV" env-default:"Dev"`
	DBName string `env:"DB_NAME" env-default:"mealdb"`
	DBUser string `env:"DB_USER" env-default:"postgres"`
	DBPass string `env:"DB_PASS" env-default:"password123"`
	DBHost string `env:"DB_HOST" env-default:"localhost"`
	DBPort string `env:"DB_PORT" env-default:"5433"`
}

func MustLoad() (*Config, error) {

	var cfg Config

	err := cleanenv.ReadEnv(&cfg)

	return &cfg, err
}
