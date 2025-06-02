package config

/*
	Service specific configuration. Extends the default configuration
*/

import (
	"context"
	"log"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/katesclau/letisgo/internal/config"
)

type BackendConfig struct {
	DynamodDBClient *dynamodb.Client
	config.Config
}

func GetConfig(ctx context.Context) *BackendConfig {

	// AWS Config
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("error loading AWS config: %v", err)
	}

	dynClient := dynamodb.NewFromConfig(awsCfg)

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		log.Fatalf("error loading default config: %v", err)
	}

	return &BackendConfig{
		Config:          *cfg,
		DynamodDBClient: dynClient,
	}
}

func (c *BackendConfig) GetConfig() *config.Config {
	return &c.Config
}
