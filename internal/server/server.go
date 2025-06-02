package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/katesclau/letisgo/internal/server/authorization"
	"github.com/katesclau/letisgo/internal/server/endpoints"
	"github.com/katesclau/letisgo/internal/server/middlewares"
	"github.com/katesclau/letisgo/internal/server/session"
)

const keyServerAddr = "mnesis-server-addr"

type Server interface {
	Listen()
}

type ServerConfig struct {
	API       endpoints.Endpoints
	RedisUrl  string
	JWTSecret string
	Port      string
}

type DefaultServer struct {
	handler http.Handler
	port    string
	name    string
	cancel  func()
	ctx     context.Context
}

func New(cfg ServerConfig) Server {
	ctx, cancel := getContext(cfg.API)

	// Create sessionManager
	authStore := authorization.New(endpoints.GetAuthorizationRoutes(cfg.API.Routes))
	sessionManager := session.New(session.RedisSessionManagerConfig{
		JWTSecret: []byte(cfg.JWTSecret),
		RedisUrl:  cfg.RedisUrl,
	}, authStore)

	wrappedMux := middlewares.ApplyMiddlewares(
		http.HandlerFunc(cfg.API.Mux.ServeHTTP),
		middlewares.LoggingMiddleware{},
		middlewares.AuthorizationMiddleware{
			SessionManager: sessionManager,
		},
		middlewares.TracingMiddleware{},
		middlewares.MetricsMiddleware{},
	)

	return &DefaultServer{
		handler: wrappedMux,
		port:    cfg.Port,
		name:    cfg.API.Name,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (s *DefaultServer) Listen() {

	instance := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: s.handler,
		BaseContext: func(listener net.Listener) context.Context {
			return context.WithValue(s.ctx, keyServerAddr, listener.Addr().String())
		},
	}

	go func() {
		if err := instance.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
			s.cancel()
		}
		log.Printf("Server started on port %s", s.port)
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer shutdownCancel()
	if err := instance.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
