package infra

import (
	"car-pooling-service/internal/car_pooling"
	dto "car-pooling-service/internal/domain"
	"encoding/json"
	"net/http"
)

type HTTPHandler struct {
	app *app.Application
}

func NewHTTPHandler(application *app.Application) *HTTPHandler {
	return &HTTPHandler{app: application}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/status":
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.handleStatus(w, r)

	case "/cars":
		if r.Method != http.MethodPut {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.handleAddCars(w, r)
	case "/journey":
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.handleAddJourney(w, r)
	case "/locate":
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.handleLocateJourney(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
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

	err := json.NewDecoder(r.Body).Decode(&addCarsRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, carRequest := range addCarsRequest {
		if !carRequest.IsValid() {
			http.Error(w, "Invalid car request data", http.StatusBadRequest)
			return
		}
	}

	err = h.app.Commands.AddCar.Handle(addCarsRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) handleAddJourney(w http.ResponseWriter, r *http.Request) {
	var addJourneyRequest dto.AddJourneyRequest

	err := json.NewDecoder(r.Body).Decode(&addJourneyRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !addJourneyRequest.IsValid() {
		http.Error(w, "Invalid journey request data", http.StatusBadRequest)
		return
	}

	err = h.app.Commands.AddJourney.Handle(addJourneyRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) handleLocateJourney(w http.ResponseWriter, r *http.Request) {
	var locateJourneyRequest dto.LocateJourneyRequest

	if err := locateJourneyRequest.Validate(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	car, err := h.app.Queries.LocateJourney.Handle(locateJourneyRequest.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// No car assigned
	if car == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(car)
}
