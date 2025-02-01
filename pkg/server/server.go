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

	"mnesis.com/pkg/server/authorization"
	"mnesis.com/pkg/server/endpoints"
	"mnesis.com/pkg/server/middlewares"
)

const keyServerAddr = "mnesis-server-addr"

type Server struct {
	handler http.Handler
	port    string
	name    string
	cancel  func()
	ctx     context.Context
}

func NewAPIDefinition(name, description string, version string, routes *endpoints.APIRoutes) *endpoints.APIDefinition {
	mux := http.ServeMux{}

	for path, endpoint := range *routes {
		mux.HandleFunc(path, endpoint.Handler)
	}

	return &endpoints.APIDefinition{
		Name:        name,
		Description: description,
		Version:     version,
		Mux:         &mux,
		Routes:      routes,
	}
}

func getAuthorizationRoutes(routes *endpoints.APIRoutes) map[string]authorization.AuthorizationRole {
	authorizationRoutes := make(map[string]authorization.AuthorizationRole)
	for path, endpoint := range *routes {
		authorizationRoutes[path] = endpoint.AuthorizationRole
	}
	return authorizationRoutes
}

func getContext(api endpoints.APIDefinition) (context.Context, context.CancelFunc) {
	ctx := context.WithValue(
		context.WithValue(
			context.WithValue(
				context.Background(),
				"api_name", api.Name,
			),
			"api_version", api.Version,
		),
		"api_description",
		api.Description,
	)
	return context.WithCancel(ctx)
}

func NewServer(api endpoints.APIDefinition, port string) *Server {
	ctx, cancel := getContext(api)

	wrappedMux := middlewares.ApplyMiddlewares(
		http.HandlerFunc(api.Mux.ServeHTTP),
		middlewares.LoggingMiddleware{},
		middlewares.AuthorizationMiddleware{
			Options: middlewares.AuthorizationMiddlewareOptions{
				AuthorizedRoutes: getAuthorizationRoutes(api.Routes),
			},
		},
		middlewares.TracingMiddleware{},
		middlewares.MetricsMiddleware{},
	)

	return &Server{
		handler: wrappedMux,
		port:    port,
		name:    api.Name,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (s *Server) Listen() {

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
