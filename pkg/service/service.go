package service

/*
Service is the main service for the API
It defines the configuration and routes for the API
It also defines the Event listeners and handlers
And starts the server
*/

import (
	"github.com/go-session/redis/v3"
	gosession "github.com/go-session/session/v3"
	"mnesis.com/pkg/config"
	"mnesis.com/pkg/server"
	"mnesis.com/pkg/server/endpoints"
	"mnesis.com/pkg/server/session"
)

type Service struct {
	Config         *config.Config
	Routes         *endpoints.APIRoutes
	SessionManager *session.SessionManager
}

func NewService(cfg *config.Config, routes *endpoints.APIRoutes) *Service {

	gosession.InitManager(
		gosession.SetStore(redis.NewRedisStore(&redis.Options{
			Addr: cfg.RedisHost,
			DB:   15,
		})),
	)

	// Create a new service
	return &Service{
		Config:         cfg,
		Routes:         routes,
		SessionManager: session.New(cfg),
	}
}

func (s *Service) Start() {
	// Start the service

	// Init Queue listener
	// TODO Implement PubSub for Async processes

	// Init HTTP Server
	api := server.NewAPIDefinition(*s.Config, s.Routes, s.SessionManager)
	server := server.NewServer(*api, s.Config.Port)
	// Start the server, and handle shutdown
	server.Listen()
}
