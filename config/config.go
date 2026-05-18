package config

import "os"

type Config struct {
	Port   string
	AppEnv string
	DBName string
	DBUser string
	DBPass string
	DBHost string
	DBPort string
}

func LoadConfig() Config {
	return Config{
		Port:   getEnv("PORT", "8085"),
		AppEnv: getEnv("APP_ENV", "Devlopment"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5433"),
		DBName: getEnv("DB_NAME", "mealdb"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPass: getEnv("DB_PASS", "password"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}
	return value
}
