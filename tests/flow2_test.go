package tests

import (
	"bytes"
	"car-pooling-service/internal"
	http2 "car-pooling-service/internal/port/http"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func buildRequest(method, url, contentType string, body []byte) (*http.Request, error) {
	reader := bytes.NewReader(body)
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return request, nil
}

type AddCarRequest struct {
	ID    int `json:"id"`
	Seats int `json:"seats"`
}

type EnqueueJourneyRequest struct {
	ID     int `json:"id"`
	People int `json:"people"`
}

func TestAddCarsConcurrently(t *testing.T) {
	// Initialize application
	application := internal.InitializeApp()

	// Create HTTP handler and server
	httpHandler := http2.NewHTTPHandler(application)
	server := httptest.NewServer(httpHandler)
	defer server.Close()

	serverURL := server.URL
	numberOfCars := 20

	// Llama a la función para enviar peticiones concurrentes
	sendConcurrentCarRequests(serverURL, numberOfCars, t, server)

	startJourneyID := 1    // ID inicial para los viajes
	numberOfJourneys := 30 // Número de viajes a generar

	// Llama a la función para enviar peticiones concurrentes
	sendConcurrentJourneyRequests(server, startJourneyID, numberOfJourneys, t)

	var wg sync.WaitGroup
	// Comprobaciones para 10 IDs de viaje aleatorios
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			randomJourneyID := rand.Intn(numberOfJourneys) + 1
			verifyJourneyStatuses(server, randomJourneyID, http.StatusNoContent, t)
		}()
	}

	// Comprobación para un ID de viaje que sabemos que no existe
	wg.Add(1)
	go func() {
		defer wg.Done()
		nonExistentJourneyID := 9999 // Un ID que sabemos que no existe
		verifyJourneyStatuses(server, nonExistentJourneyID, http.StatusNotFound, t)
	}()

	wg.Wait() // Esperar a que todas las Go rutinas terminen

	var wg2 sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			randomJourneyDropoffID := rand.Intn(numberOfJourneys) + 1
			dropOffJourney(server, randomJourneyDropoffID, http.StatusNoContent, t)
		}()
	}

	// Comprobación para un ID de viaje que sabemos que no existe
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		nonExistentJourneyID := 9999 // Un ID que sabemos que no existe
		dropOffJourney(server, nonExistentJourneyID, http.StatusNotFound, t)
	}()

	wg2.Wait() // Esperar a que todas las Go rutinas terminen
}

func sendConcurrentCarRequests(serverURL string, numberOfCars int, t *testing.T, server *httptest.Server) {
	var wg sync.WaitGroup

	for i := 0; i < numberOfCars; i++ {
		wg.Add(1)
		go func(carID int) {
			defer wg.Done()

			seats := rand.Intn(3) + 4 // random 4, 5 or 6

			car := AddCarRequest{
				ID:    carID,
				Seats: seats,
			}
			carJson, _ := json.Marshal([]AddCarRequest{car})

			request, err := buildRequest(http.MethodPut, serverURL+"/cars", "application/json", carJson)
			if err != nil {
				t.Errorf("Failed to create PUT request: %v", err)
				return
			}

			response, err := server.Client().Do(request)
			if err != nil {
				t.Errorf("Failed to send PUT request: %v", err)
				return
			}

			if response.StatusCode != http.StatusOK {
				t.Errorf("Failed to add car with ID %d: response status %d", carID, response.StatusCode)
			}
		}(i)
	}

	wg.Wait() // Esperar a que todas las Go rutinas terminen
}

func sendConcurrentJourneyRequests(server *httptest.Server, startJourneyID, numberOfJourneys int, t *testing.T) {
	var wg sync.WaitGroup
	serverURL := server.URL

	for i := 0; i < numberOfJourneys; i++ {
		wg.Add(1)
		go func(journeyID int) {
			defer wg.Done()

			// Generar un número aleatorio de personas entre 1 y 6
			people := rand.Intn(6) + 1

			journey := EnqueueJourneyRequest{
				ID:     journeyID,
				People: people,
			}
			journeyJson, _ := json.Marshal(journey)

			request, err := http.NewRequest(http.MethodPost, serverURL+"/journey", bytes.NewBuffer(journeyJson))
			request.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Errorf("Failed to create POST request: %v", err)
				return
			}

			response, err := server.Client().Do(request)
			if err != nil {
				t.Errorf("Failed to send POST request: %v", err)
				return
			}

			if response.StatusCode != http.StatusAccepted {
				t.Errorf("Failed to enqueue journey with ID %d: response status %d", journeyID, response.StatusCode)
			}
			//log.Printf("Request Add Journey: %d with %d\n", journey.ID, journey.People)
		}(startJourneyID + i)
	}

	wg.Wait() // Esperar a que todas las Go rutinas terminen
}

func verifyJourneyStatuses(server *httptest.Server, journeyID int, expectedStatusCode int, t *testing.T) {

	request, err := buildRequest(http.MethodPost, server.URL+"/locate", "application/x-www-form-urlencoded", []byte(fmt.Sprintf("ID=%d", journeyID)))
	if err != nil {
		t.Errorf("Failed to create POST request: %v", err)
		return
	}

	response, err := server.Client().Do(request)
	if err != nil {
		t.Errorf("Failed to send POST request: %v", err)
		return
	}

	if response.StatusCode != expectedStatusCode {
		t.Errorf("Unexpected status %d for journey %d, expected %d", response.StatusCode, journeyID, expectedStatusCode)
	}

}

func dropOffJourney(server *httptest.Server, journeyID int, expectedStatusCode int, t *testing.T) {
	request, err := buildRequest(http.MethodPost, server.URL+"/dropoff", "application/x-www-form-urlencoded", []byte(fmt.Sprintf("ID=%d", journeyID)))
	if err != nil {
		t.Errorf("Failed to create POST request: %v", err)
		return
	}

	response, err := server.Client().Do(request)
	if err != nil {
		t.Errorf("Failed to send POST request: %v", err)
		return
	}

	if response.StatusCode != expectedStatusCode {
		t.Errorf("Unexpected dropoff %d for journey %d, expected %d", response.StatusCode, journeyID, expectedStatusCode)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
