package main

import (
	"car-pooling-service/internal/car_pooling"
	"car-pooling-service/internal/car_pooling/infra"
	"log"
	"net/http"
	"time"
)

func main() {
	// Init App and Configs
	application := app.InitializeApp()

	// Init Coroutine to check Journey/Car assigning
	go func() {
		for {
			application.Commands.AssignCarsToJourneys.Handle()
			time.Sleep(time.Second * 5)
		}
	}()

	httpHandler := infra.NewHTTPHandler(application)

	// Start HTTP Server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
