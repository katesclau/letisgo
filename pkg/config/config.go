package config

import (
	"context"
	"log"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `env:"ENVIRONMENT"`
	Port           string `env:"PORT"`
	APIName        string `env:"API_NAME"`
	APIVersion     string `env:"API_VERSION"`
	APIDescription string `env:"API_DESCRIPTION"`
}

func GetConfig(ctx context.Context) (*Config, error) {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err = env.UnmarshalFromEnviron(cfg)
	if err != nil {
		return nil, err
	}

	fillDefaultValues(cfg)

	return cfg, nil
}

func fillDefaultValues(cfg *Config) {
	cfg.Env = "dev"
	cfg.Port = "8080"
	cfg.APIName = "letisgo"
	cfg.APIVersion = "0.0.0"
	cfg.APIDescription = "This is an service boilerplate."
}
