package http

import (
	"car-pooling-service/internal"
	"car-pooling-service/internal/port/http/dto"
	"encoding/json"
	"log"
	"net/http"
)

type HTTPHandler struct {
	app *internal.Application
}

func NewHTTPHandler(app *internal.Application) *HTTPHandler {
	return &HTTPHandler{app: app}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/status":
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.handleStatus(w, r)

	case "/cars":
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.handleAddCars(w, r)
	case "/journey":
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.handleEnqueueJourney(w, r)
	case "/locate":
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.handleLocateJourney(w, r)
	case "/dropoff":
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.handleDropoff(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

// HTTP Ports.
// Each port handler includes client request validation
func (h *HTTPHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	h.app.Queries.Status.Handle()
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) handleAddCars(w http.ResponseWriter, r *http.Request) {
	var addCarsRequest []dto.AddCarRequest

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&addCarsRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, carRequest := range addCarsRequest {
		if !carRequest.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	err = h.app.Commands.AddCar.Handle(addCarsRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) handleEnqueueJourney(w http.ResponseWriter, r *http.Request) {
	var enqueueJourneyRequest dto.EnqueueJourneyRequest

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&enqueueJourneyRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !enqueueJourneyRequest.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.app.Commands.EnqueueJourney.Handle(enqueueJourneyRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *HTTPHandler) handleLocateJourney(w http.ResponseWriter, r *http.Request) {
	var locateJourneyRequest dto.LocateJourneyRequest

	if err := locateJourneyRequest.Validate(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	car, err := h.app.Queries.LocateJourney.Handle(locateJourneyRequest.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if car == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	carResponse := dto.CarResponse{
		ID:    car.ID,
		Seats: car.Seats,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(carResponse); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (h *HTTPHandler) handleDropoff(w http.ResponseWriter, r *http.Request) {
	var dropoffRequest dto.DropoffRequest

	if err := dropoffRequest.Validate(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.app.Commands.Dropoff.Handle(dropoffRequest)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
