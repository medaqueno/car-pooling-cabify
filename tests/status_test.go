package tests

import (
	"car-pooling-service/internal"
	"car-pooling-service/internal/infrastructure/port"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatus(t *testing.T) {
	// Arrange
	application := internal.InitializeApp()
	httpHandler := port.NewHTTPHandler(application)
	server := httptest.NewServer(httpHandler)
	defer server.Close()

	// Act
	resp, err := http.Get(server.URL + "/status")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}