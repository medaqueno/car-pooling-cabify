package command

import (
	"car-pooling-service/internal/car_pooling/repository"
	"car-pooling-service/internal/domain"
	"fmt"
)

type AddJourneyHandler struct {
	repo repository.JourneyRepository
}

func NewAddJourneyHandler(repo repository.JourneyRepository) *AddJourneyHandler {
	return &AddJourneyHandler{
		repo: repo,
	}
}

func (h *AddJourneyHandler) Handle(addJourneyRequest dto.AddJourneyRequest) error {
	journey := dto.NewJourney(addJourneyRequest.ID, addJourneyRequest.People)
	err := h.repo.AddJourney(journey)
	if err != nil {
		return fmt.Errorf("Error adding journey: %v\n", err)
	}

	// Debug
	h.repo.LogAllJourneys()

	return nil
}
