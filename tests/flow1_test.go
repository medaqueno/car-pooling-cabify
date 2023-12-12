package tests

import (
	"bytes"
	"car-pooling-service/internal"
	"car-pooling-service/internal/infrastructure/port"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Car struct {
	ID    int `json:"id"`
	Seats int `json:"seats"`
}

type Journey struct {
	ID     int `json:"id"`
	People int `json:"people"`
}

func buildRequest(method, url, contentType string, body []byte) (*http.Request, error) {
	reader := bytes.NewReader(body)
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return request, nil
}

type JourneyVerification struct {
	ID                 int `json:"id"`
	ExpectedStatusCode int `json:"expected_status_code"`
}

func verifyJourneyStatuses(t *testing.T, server *httptest.Server, journeys []JourneyVerification) {
	for _, journey := range journeys {
		request, err := buildRequest(http.MethodPost, server.URL+"/locate", "application/x-www-form-urlencoded", []byte(fmt.Sprintf("ID=%d", journey.ID)))
		if err != nil {
			t.Errorf("Failed to create POST request: %v", err)
			return
		}

		response, err := server.Client().Do(request)
		if err != nil {
			t.Errorf("Failed to send POST request: %v", err)
			return
		}

		if response.StatusCode != journey.ExpectedStatusCode {
			t.Errorf("Unexpected status %d for journey %d, expected %d", response.StatusCode, journey.ID, journey.ExpectedStatusCode)
		}
	}
}

func TestCarAssignmentFlow(t *testing.T) {
	// Initialize application
	application := internal.InitializeApp()

	// Create HTTP handler and server
	httpHandler := port.NewHTTPHandler(application)
	server := httptest.NewServer(httpHandler)
	defer server.Close()

	// Define expected journey statuses for various scenarios
	initialJourneyVerification := []JourneyVerification{
		{1, 200},
		{2, 204},
		/*{3, 200},
		{4, 200},
		{5, 200},*/
		{6, 204},
		{7, 204},
	}

	additionalCarVerification := []JourneyVerification{
		{6, 204},
		{7, 204},
	}

	remainingCarVerification := []JourneyVerification{
		{7, 204},
	}

	t.Run("Initial car addition and journey assignment", func(t *testing.T) {
		// Create initial cars
		initialCars := []Car{
			{1, 4},
			/*{2, 5},
			{3, 6},
			{4, 6},
			{5, 6},*/
		}

		// Create journeys
		initialJourneys := []Journey{
			{1, 4},
			{2, 3},
			/*{3, 6},
			{4, 5},
			{5, 4},*/
			{6, 2},
			{7, 6},
		}

		// Add initial cars
		data, err := json.Marshal(initialCars)
		if err != nil {
			t.Errorf("Failed to marshal car data: %v", err)
			return
		}

		request, err := buildRequest(http.MethodPut, server.URL+"/cars", "application/json", data)
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
			t.Errorf("Unexpected status code: %d", response.StatusCode)
		}

		// Enqueue journeys
		for _, journey := range initialJourneys {
			data, err := json.Marshal(journey)
			if err != nil {
				t.Errorf("Failed to marshal journey data: %v", err)
				return
			}

			request, err := buildRequest(http.MethodPost, server.URL+"/journey", "application/json", data)
			if err != nil {
				t.Errorf("Failed to create POST request: %v", err)
				return
			}

			response, err := server.Client().Do(request)
			if err != nil {
				t.Errorf("Failed to send POST request: %v", err)
				return
			}

			if response.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code: %d", response.StatusCode)
			}
		}

		// Assign cars
		application.Services.CarAssigner.AssignCarsToJourneys()

		// Verify journey statuses
		verifyJourneyStatuses(t, server, initialJourneyVerification)
	})

	t.Run("Adding additional cars and re-assignment", func(t *testing.T) {
		// Add additional cars
		additionalCars := []Car{
			//{6, 4},
			//	{7, 5},
		}

		// Marshal car data
		data, err := json.Marshal(additionalCars)
		if err != nil {
			t.Errorf("Failed to marshal car data: %v", err)
			return
		}

		// Build PUT request for adding additional cars
		request, err := buildRequest(http.MethodPut, server.URL+"/cars", "application/json", data)
		if err != nil {
			t.Errorf("Failed to create PUT request: %v", err)
			return
		}

		// Send PUT request
		response, err := server.Client().Do(request)
		if err != nil {
			t.Errorf("Failed to send PUT request: %v", err)
			return
		}

		// Check PUT request status code
		if response.StatusCode != http.StatusOK {
			t.Errorf("Unexpected status code: %d", response.StatusCode)
		}

		// Assign cars
		application.Services.CarAssigner.AssignCarsToJourneys()

		// Verify journey statuses
		verifyJourneyStatuses(t, server, additionalCarVerification)
	})

	t.Run("Adding cars and assigning remaining journeys", func(t *testing.T) {
		// Add remaining cars
		remainingCars := []Car{
			//	{8, 6}
		}

		// Marshal car data
		data, err := json.Marshal(remainingCars)
		if err != nil {
			t.Errorf("Failed to marshal car data: %v", err)
			return
		}

		// Build PUT request for adding remaining cars
		request, err := buildRequest(http.MethodPut, server.URL+"/cars", "application/json", data)
		if err != nil {
			t.Errorf("Failed to create PUT request: %v", err)
			return
		}

		// Send PUT request
		response, err := server.Client().Do(request)
		if err != nil {
			t.Errorf("Failed to send PUT request: %v", err)
			return
		}

		// Check PUT request status code
		if response.StatusCode != http.StatusOK {
			t.Errorf("Unexpected status code: %d", response.StatusCode)
		}

		// Assign cars
		application.Services.CarAssigner.AssignCarsToJourneys()

		// Verify journey statuses
		verifyJourneyStatuses(t, server, remainingCarVerification)
	})

	t.Run("Drop-off and assigning new journey", func(t *testing.T) {
		// Define new journey data
		newJourney := Journey{
			ID:     8,
			People: 3,
		}

		// Drop-off existing journey
		dropOffJourneyID := 1 // Replace with actual journey ID
		request, err := buildRequest(http.MethodPost, server.URL+"/dropoff", "application/x-www-form-urlencoded", []byte(fmt.Sprintf("ID=%d", dropOffJourneyID)))
		if err != nil {
			t.Errorf("Failed to create POST request: %v", err)
			return
		}

		response, err := server.Client().Do(request)
		if err != nil {
			t.Errorf("Failed to send POST request: %v", err)
			return
		}

		if response.StatusCode != http.StatusNoContent {
			t.Errorf("Unexpected status code: %d for drop-off", response.StatusCode)
		}

		// Add new journey
		data, err := json.Marshal(newJourney)
		if err != nil {
			t.Errorf("Failed to marshal journey data: %v", err)
			return
		}

		request, err = buildRequest(http.MethodPost, server.URL+"/journey", "application/json", data)
		if err != nil {
			t.Errorf("Failed to create POST request: %v", err)
			return
		}

		response, err = server.Client().Do(request)
		if err != nil {
			t.Errorf("Failed to send POST request: %v", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Unexpected status code: %d for new journey", response.StatusCode)
		}

		// Assign cars
		application.Services.CarAssigner.AssignCarsToJourneys()
		/*
			updatedVerification := []JourneyVerification{
				{1, 404},
				{2, 200},
				{3, 200},
				{4, 200},
				{5, 200},
				{6, 200},
				{7, 200},
				{newJourney.ID, 200},
			}

			// Verify journey statuses
			verifyJourneyStatuses(t, server, updatedVerification)
		*/
	})

}
