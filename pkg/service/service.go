package service

/*
Service is the main service for the API
It defines the configuration and routes for the API
It also defines the Event listeners and handlers
And starts the server
*/

import (
	"mnesis.com/pkg/config"
	"mnesis.com/pkg/server"
	"mnesis.com/pkg/server/endpoints"
)

type Service struct {
	Config *config.Config
	Routes *endpoints.APIRoutes
}

func NewService(cfg *config.Config, routes *endpoints.APIRoutes) *Service {
	// Create a new service
	return &Service{
		Config: cfg,
		Routes: routes,
	}
}

func (s *Service) Start() {
	// Start the service

	// Init Queue listener
	// TODO Implement PubSub for Async processes

	// Init HTTP Server
	api := server.NewAPIDefinition(s.Config.APIName, s.Config.APIDescription, s.Config.APIVersion, s.Routes)
	server := server.NewServer(*api, s.Config.Port)
	// Start the server, and handle shutdown
	server.Listen()
}
