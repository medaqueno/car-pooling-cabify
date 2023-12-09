package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Working...")
	})

	// Start the server
	log.Println("Init server in port: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
