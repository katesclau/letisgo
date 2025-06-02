/**
This is an API server for Price definition of medical and dental procedures.
*/

package main

import (
	"context"

	"github.com/katesclau/letisgo/backend/routes"
	"github.com/katesclau/letisgo/internal/config"
	"github.com/katesclau/letisgo/internal/service"
	log "github.com/sirupsen/logrus"
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
