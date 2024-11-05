package endpoints

import (
	"encoding/json"
	"net/http"
)

var Health = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		apiName := r.Context().Value("api_name")
		apiVersion := r.Context().Value("api_version")
		apiDescription := r.Context().Value("api_description")

		if apiName == nil || apiVersion == nil {
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
			http.Error(w, "Failed to generate health response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})
