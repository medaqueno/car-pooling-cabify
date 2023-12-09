package main

import (
	"car-pooling-service/internal/car_pooling"
	"car-pooling-service/internal/car_pooling/infra"
	"log"
	"net/http"
)

func main() {
	// Init App and Configs
	application := app.InitializeApp()
	httpHandler := infra.NewHTTPHandler(application)

	// Start HTTP Server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
