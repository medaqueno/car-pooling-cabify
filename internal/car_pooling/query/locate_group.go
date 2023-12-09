package query

import (
	"car-pooling-service/internal/car_pooling/repository"
	dto "car-pooling-service/internal/domain"
	"fmt"
)

type LocateCarByJourneyHandler struct {
	journeyRepo repository.JourneyRepository
	carRepo     repository.CarRepository
}

func NewLocateJourneyHandler(journeyRepo repository.JourneyRepository, carRepo repository.CarRepository) *LocateCarByJourneyHandler {
	return &LocateCarByJourneyHandler{
		journeyRepo: journeyRepo,
		carRepo:     carRepo,
	}
}

func (h *LocateCarByJourneyHandler) Handle(groupID int) (*dto.Car, error) {
	journey, err := h.journeyRepo.FindJourneyByID(groupID)
	// No Journey
	if err != nil {
		fmt.Printf("No Journey Found\n")
		return nil, err
	}

	// No Car Assigned
	if journey.CarId == nil {
		fmt.Printf("No Car Assigned\n")
		return nil, nil
	}

	car, err := h.carRepo.FindCarByID(*journey.CarId)
	// Car does not exist
	if err != nil {
		fmt.Printf("Car does not exist\n")
		return nil, err
	}

	return car, nil
}
