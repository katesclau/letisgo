package server

import (
	"net/http"

	"mnesis.com/pkg/config"
	"mnesis.com/pkg/server/endpoints"
)

func NewAPIDefinition(cfg config.Config, routes *endpoints.Routes) endpoints.Endpoints {
	mux := http.ServeMux{}

	for path, endpoint := range *routes {
		mux.HandleFunc(path, endpoint.Handler)
	}

	return endpoints.Endpoints{
		Name:        cfg.ServiceName,
		Description: cfg.ServiceDescription,
		Version:     cfg.ServiceVersion,
		Mux:         &mux,
		Routes:      routes,
	}
}
