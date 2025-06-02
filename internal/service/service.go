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
	"github.com/katesclau/letisgo/internal/config"
	"github.com/katesclau/letisgo/internal/server"
	"github.com/katesclau/letisgo/internal/server/endpoints"
)

type Service struct {
	Config *config.Config
	Routes *endpoints.Routes
}

func NewService(cfg *config.Config, routes *endpoints.Routes) *Service {

	gosession.InitManager(
		gosession.SetStore(redis.NewRedisStore(&redis.Options{
			Addr: cfg.RedisHost,
			DB:   15,
		})),
	)

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
	api := server.NewAPIDefinition(*s.Config, s.Routes)

	server := server.New(server.ServerConfig{
		API:       api,
		RedisUrl:  s.Config.RedisHost,
		JWTSecret: s.Config.JWTSecret,
		Port:      s.Config.Port,
	})
	// Start the server, and handle shutdown
	server.Listen()
}
