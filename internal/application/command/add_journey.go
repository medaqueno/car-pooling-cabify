package command

import (
	"car-pooling-service/internal/domain/model"
	"car-pooling-service/internal/domain/repository"
	"car-pooling-service/internal/port/http/dto"
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

	journey := model.NewJourney(enqueueJourneyRequest.ID, enqueueJourneyRequest.People)
	err := h.repo.EnqueueJourney(journey)
	if err != nil {
		return fmt.Errorf("Error enqueueing journey: %v\n", err)
	}

	return nil
}
