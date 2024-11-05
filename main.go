/**
This is an API server for Price definition of medical and dental procedures.
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"mnesis.com/pkg/server"
	"mnesis.com/pkg/server/endpoints"
)

func main() {
	fmt.Println("Starting Service...")

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	apiName := os.Getenv("API_NAME")
	apiDescription := os.Getenv("API_DESCRIPTION")
	apiVersion := os.Getenv("API_VERSION")

	// Init Queue listener
	// TODO Implement PubSub for Async processes

	// Init HTTP Server
	api := endpoints.NewAPIDefinition(apiName, apiDescription, apiVersion)
	server := server.NewServer(*api, port)
	server.Listen()
}
