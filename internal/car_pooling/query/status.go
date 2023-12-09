package query

import "log"

// StatusHandler Dependencies used by action
type StatusHandler struct {
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{
		// Init dependencies
	}
}

func (h *StatusHandler) Handle() {
	log.Println("Status OK")
}
