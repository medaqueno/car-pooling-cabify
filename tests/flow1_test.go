package tests

import (
	"bytes"
	app "car-pooling-service/internal/car_pooling"
	"car-pooling-service/internal/car_pooling/infra"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestCarAssignmentFlow(t *testing.T) {
	// Initialize application
	application := app.InitializeApp()

	// Create HTTP handler and server
	httpHandler := infra.NewHTTPHandler(application)
	server := httptest.NewServer(httpHandler)
	defer server.Close()

	// Register Cars
	// Create car objects with different seat capacities
	cars := []struct {
		ID    int `json:"id"`
		Seats int `json:"seats"`
	}{
		{1, 4},
		{2, 5},
		{3, 6},
		{4, 6},
		{5, 6},
	}

	requestAddCars(t, server, cars)

	// Enqueue journeys
	// Create journey objects with varying group sizes
	journeys := []struct {
		ID     int `json:"id"`
		People int `json:"people"`
	}{
		{1, 4},
		{2, 3},
		{3, 6},
		{4, 5},
		{5, 4},
		{6, 2},
		{7, 6},
	}

	// Enqueue journeys using POST request to /journey endpoint
	for _, journey := range journeys {
		requestEnqueueJourneys(t, server, journey)
	}

	application.Commands.AssignCarsToJourneys.Handle()

	// Check Journeys
	// Check locations of all journeys using POST request to /locate endpoint
	expected := map[int]int{
		1: 200, 2: 200, 3: 200,
		4: 200, 5: 200, 6: 204,
		7: 204,
	}
	// Check Journeys
	// Check locations of all journeys using POST request to /locate endpoint
	var wg sync.WaitGroup
	for journeyID, statusCode := range expected {
		wg.Add(1)
		go func(journeyID int, statusCode int) {
			defer wg.Done()
			fmt.Printf("Processing groupID %v\n", journeyID)
			verifyJourneyStatus(t, server, journeyID, statusCode)
		}(journeyID, statusCode)
	}
	wg.Wait()

	// Add More Cars
	cars2 := []struct {
		ID    int `json:"id"`
		Seats int `json:"seats"`
	}{
		{6, 4},
		{7, 5},
	}

	requestAddCars(t, server, cars2)
	application.Commands.AssignCarsToJourneys.Handle()

	expected2 := map[int]int{
		6: 200,
		7: 204,
	}
	var wg2 sync.WaitGroup
	for journeyID, statusCode := range expected2 {
		wg2.Add(1)
		go func(journeyID int, statusCode int) {
			defer wg2.Done()
			fmt.Printf("Processing groupID %v\n", journeyID)
			verifyJourneyStatus(t, server, journeyID, statusCode)
		}(journeyID, statusCode)
	}
	wg2.Wait()

	// Add More Cars
	cars3 := []struct {
		ID    int `json:"id"`
		Seats int `json:"seats"`
	}{
		{8, 6},
	}

	requestAddCars(t, server, cars3)
	application.Commands.AssignCarsToJourneys.Handle()
	expected3 := map[int]int{
		7: 200,
	}
	var wg3 sync.WaitGroup
	for journeyID, statusCode := range expected3 {
		wg3.Add(1)
		go func(journeyID int, statusCode int) {
			defer wg3.Done()
			fmt.Printf("Processing groupID %v\n", journeyID)
			verifyJourneyStatus(t, server, journeyID, statusCode)
		}(journeyID, statusCode)
	}
	wg3.Wait()

}

func requestAddCars(t *testing.T, server *httptest.Server, cars []struct {
	ID    int `json:"id"`
	Seats int `json:"seats"`
}) {
	data, err := json.Marshal(cars)
	if err != nil {
		t.Errorf("Failed to marshal car data: %v", err)
		return
	}

	request, err := http.NewRequest(http.MethodPut, server.URL+"/cars", bytes.NewReader(data))
	if err != nil {
		t.Errorf("Failed to create PUT request: %v", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	responseAddCars, err := server.Client().Do(request)
	if err != nil {
		t.Errorf("Failed to send PUT request: %v", err)
		return
	}

	if responseAddCars.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", responseAddCars.StatusCode)
	}

}

func requestEnqueueJourneys(t *testing.T, server *httptest.Server, journey struct {
	ID     int `json:"id"`
	People int `json:"people"`
}) {
	data, err := json.Marshal(journey)
	if err != nil {
		t.Errorf("Failed to marshal journey data: %v", err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, server.URL+"/journey", bytes.NewReader(data))
	if err != nil {
		t.Errorf("Failed to create POST request: %v", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := server.Client().Do(request)
	if err != nil {
		t.Errorf("Failed to send POST request: %v", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Unexpected status code: %d", response.StatusCode)
	}

}

func verifyJourneyStatus(t *testing.T, server *httptest.Server, id int, statusCode int) {
	request, err := http.NewRequest(http.MethodPost, server.URL+"/locate", strings.NewReader(fmt.Sprintf("ID=%d", id)))
	if err != nil {
		t.Errorf("Failed to create POST request: %v", err)
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := server.Client().Do(request)
	if err != nil {
		t.Errorf("Failed to send POST request: %v", err)
		return
	}

	if response.StatusCode != statusCode {
		t.Errorf("Unexpected status %d for journey %d, expected %d", response.StatusCode, id, statusCode)
	}
}
