package main

import (
	"car-pooling-service/internal"
	"car-pooling-service/internal/infrastructure/config"
	httpPort "car-pooling-service/internal/port/http"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Init App
	application := internal.InitializeApp()

	// Init Coroutine to check Journey/Car assigning
	go func() {
		for {
			application.Services.CarAssigner.Handle()
		}
	}()

	httpHandler := httpPort.NewHTTPHandler(application)

	// Start HTTP Server
	log.Printf("Starting server on %s in %s environment", cfg.ServerPort, cfg.AppEnvironment)
	if err := http.ListenAndServe(":"+cfg.ServerPort, httpHandler); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
