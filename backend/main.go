/**
This is an API server for Price definition of medical and dental procedures.
*/

package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"mnesis.com/backend/routes"
	"mnesis.com/pkg/config"
	"mnesis.com/pkg/service"
)

func main() {
	log.Info("Starting Service...")
	ctx := context.Background()

	// Load Configuration
	cfg, err := config.GetConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Define Log Level
	log.SetLevel(cfg.LogLevel)

	// Load Routes
	routes := routes.Get()

	// Init Service
	service := service.NewService(cfg, routes)
	service.Start()

}
