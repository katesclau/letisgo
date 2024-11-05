package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var Health = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		apiName := r.Context().Value("api_name")
		apiVersion := r.Context().Value("api_version")
		apiDescription := r.Context().Value("api_description")

		fmt.Printf("Request: %v", r)

		if apiName == nil || apiVersion == nil {
			fmt.Printf("Missing API information in context")
			http.Error(w, "Missing API information in context", http.StatusBadRequest)
			return
		}
		apiInfo := map[string]string{
			"name":        apiName.(string),
			"description": apiDescription.(string),
			"version":     apiVersion.(string),
		}
		jsonResponse, err := json.Marshal(apiInfo)
		if err != nil {
			fmt.Printf("Failed to generate health response %v", err)
			http.Error(w, "Failed to generate health response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})
