package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppPort string `yaml:"app_port" env:"PORT"`
	DBHost  string `yaml:"db_host" env:"DB_HOST"`
	DBPort  string `yaml:"db_port" env:"DB_PORT"`
	DBUser  string `yaml:"db_user" env:"DB_USER"`
	DBPass  string `yaml:"db_pass" env:"DB_PASS"`
	DBName  string `yaml:"db_name" env:"DB_NAME"`
}

func MustLoad() (*Config, error) {

	var configPath string

	configPath = os.Getenv("CONFIG_PATH") //Check if it is available at the given location

	if configPath == "" {
		flags := flag.String("config", "", "configuratio file path")
		// this is passed in CLI in this format --namespace=
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path not found")
		}
		//Check if the file does not exists at the given location
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file does not exist: %s", configPath)
		}
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		return nil, err
	}

	// Override using ENV vars
	err = cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
