package main

import (
	"fmt"
	"net/http"
	"os"

	"rssreader/internal/pkg/http/handlers"
	"rssreader/internal/pkg/http/routing"
	"rssreader/internal/pkg/sources"
)

func main() {
	serverAddress := os.Getenv("RSS_READER_SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = ":8080"
	}

	sourceDefinition, err := sources.New()
	if err != nil {
		fmt.Printf("Error reading the source definitions: %v", err)
		return
	}

	sourcesHandler := handlers.NewSources(sourceDefinition)
	feedsHandler := handlers.NewArticles(sourceDefinition)
	mux := routing.New(sourcesHandler, feedsHandler)

	if err := http.ListenAndServe(serverAddress, mux); err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
	}
}
