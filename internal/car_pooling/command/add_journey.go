package command

import (
	"car-pooling-service/internal/car_pooling/repository"
	"car-pooling-service/internal/domain"
	"fmt"
)

type EnqueueJourneyHandler struct {
	repo repository.JourneyRepository
}

func NewEnqueueJourneyHandler(repo repository.JourneyRepository) *EnqueueJourneyHandler {
	return &EnqueueJourneyHandler{
		repo: repo,
	}
}

func (h *EnqueueJourneyHandler) Handle(enqueueJourneyRequest dto.EnqueueJourneyRequest) error {
	journey := dto.NewJourney(enqueueJourneyRequest.ID, enqueueJourneyRequest.People)
	err := h.repo.EnqueueJourney(journey)
	if err != nil {
		return fmt.Errorf("Error enqueueing journey: %v\n", err)
	}

	// Debug
	// h.repo.LogAllJourneys()

	return nil
}
