package query

import (
	"car-pooling-service/internal/domain/model"
	"car-pooling-service/internal/domain/repository"
)

type LocateCarByJourneyHandler struct {
	repo repository.CarAssignerRepository
}

func NewLocateJourneyHandler(repo repository.CarAssignerRepository) *LocateCarByJourneyHandler {
	return &LocateCarByJourneyHandler{
		repo: repo,
	}
}

func (h *LocateCarByJourneyHandler) Handle(journeyID int) (*model.Car, error) {
	return h.repo.FindCarByJourneyID(journeyID)
}
