package endpoints

import "net/http"

type APIDefinition struct {
	Name        string
	Description string
	Version     string
	Mux         *http.ServeMux
}

func NewAPIDefinition(name, description string, version string) *APIDefinition {
	mux := http.ServeMux{}

	// Declares Paths
	mux.HandleFunc("/", Root.ServeHTTP)
	mux.HandleFunc("/health", Health.ServeHTTP)

	return &APIDefinition{
		Name:        name,
		Description: description,
		Version:     version,
		Mux:         &mux,
	}
}
