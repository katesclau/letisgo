/**
This is an API server for Price definition of medical and dental procedures.
*/

package main

import (
	"context"
	"fmt"
	"log"

	"mnesis.com/backend/routes"
	"mnesis.com/pkg/config"
	"mnesis.com/pkg/service"
)

func main() {
	fmt.Println("Starting Service...")
	ctx := context.TODO()

	// Load Configuration
	cfg, err := config.GetConfig(ctx)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Load Routes
	routes := routes.Get()

	// Init Service
	service := service.NewService(cfg, routes)
	service.Start()

}
