package main

import (
	"car-pooling-service/internal"
	"car-pooling-service/internal/infrastructure/config"
	"car-pooling-service/internal/infrastructure/port"
	"log"
	"net/http"
)

func main() {

	// Init Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Init App
	application := internal.InitializeApp()

	// Init Coroutine to check Journey/Car assigning
	application.Services.CarAssigner.RunAssignmentProcess()

	httpHandler := port.NewHTTPHandler(application)

	// Start HTTP Server
	log.Printf("Starting server on %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, httpHandler); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
