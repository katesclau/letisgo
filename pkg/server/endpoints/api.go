package endpoints

import "net/http"

type APIDefinition struct {
	Name        string
	Description string
	Version     string
	Mux         *http.ServeMux
}

type APIRoutes map[string]http.HandlerFunc

func NewAPIDefinition(name, description string, version string, routes *APIRoutes) *APIDefinition {
	mux := http.ServeMux{}

	for path, handler := range *routes {
		mux.HandleFunc(path, handler)
	}

	return &APIDefinition{
		Name:        name,
		Description: description,
		Version:     version,
		Mux:         &mux,
	}
}
