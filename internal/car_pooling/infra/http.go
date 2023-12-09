package infra

import (
	"car-pooling-service/internal/car_pooling"
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
	}
}

// HTTP Ports
func (h *HTTPHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	h.app.Queries.Status.Handle()
	w.WriteHeader(http.StatusOK)
}
