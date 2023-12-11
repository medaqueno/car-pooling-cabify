package main

import (
	"car-pooling-service/internal"
	"car-pooling-service/internal/infrastructure/port"
	"log"
	"net/http"
	"time"
)

func main() {
	// Init App and Configs
	application := internal.InitializeApp()

	// Init Coroutine to check Journey/Car assigning
	go func() {
		for {
			application.Services.AssignCarsToJourneys.Handle()
			time.Sleep(time.Second * 5)
		}
	}()

	httpHandler := port.NewHTTPHandler(application)

	// Start HTTP Server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
