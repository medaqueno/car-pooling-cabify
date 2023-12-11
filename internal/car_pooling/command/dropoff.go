package command

import (
	"car-pooling-service/internal/car_pooling/repository"
	dto "car-pooling-service/internal/domain"
	"fmt"
)

type DropoffHandler struct {
	carRepo     repository.CarRepository
	journeyRepo repository.JourneyRepository
}

func NewDropoffHandler(carRepo repository.CarRepository, journeyRepo repository.JourneyRepository) *DropoffHandler {
	return &DropoffHandler{
		carRepo:     carRepo,
		journeyRepo: journeyRepo,
	}
}

func (h *DropoffHandler) Handle(dropoffRequest dto.DropoffRequest) error {
	journey, err := h.journeyRepo.FindJourneyByID(dropoffRequest.ID)
	if err != nil {
		return fmt.Errorf("error finding journey: %v", err)
	}

	if journey.CarId != nil {
		car, err := h.carRepo.FindCarByID(*journey.CarId)
		if err != nil {
			return fmt.Errorf("error finding car for journey: %v", err)
		}
		car.AvailableSeats += journey.People
	}

	err = h.journeyRepo.RemoveJourney(journey.ID)
	if err != nil {
		return fmt.Errorf("error removing journey: %v", err)
	}

	return nil
}
