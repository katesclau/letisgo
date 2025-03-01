package config

import (
	"context"
	"log"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Env                string `env:"ENVIRONMENT"`
	Port               string `env:"PORT"`
	ServiceName        string `env:"SERVICE_NAME"`
	ServiceVersion     string `env:"SERVICE_VERSION"`
	ServiceDescription string `env:"SERVICE_DESCRIPTION"`
	RedisHost          string `env:"REDIS_HOST"`

	JWTSecret string `env:"JWT_SECRET,default:jwt-secret"`

	LogLevelString string `env:"LOG_LEVEL"`
	LogLevel       logrus.Level
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
	cfg.Env = firstNonEmpty(cfg.Env, "dev")
	cfg.Port = firstNonEmpty(cfg.Port, "8080")
	cfg.ServiceName = firstNonEmpty(cfg.ServiceName, "letisgo")
	cfg.ServiceVersion = firstNonEmpty(cfg.ServiceVersion, "0.0.0")
	cfg.ServiceDescription = firstNonEmpty(cfg.ServiceDescription, "This is a service boilerplate.")

	setLogLevel(cfg)
}

func setLogLevel(cfg *Config) {
	level, err := logrus.ParseLevel(cfg.LogLevelString)
	if err != nil {
		cfg.LogLevel = logrus.InfoLevel
	} else {
		cfg.LogLevel = level
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
