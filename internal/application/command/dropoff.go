package command

import (
	"car-pooling-service/internal/application/service"
	"car-pooling-service/internal/domain/model"
	"car-pooling-service/internal/domain/repository"
	"fmt"
)

type DropoffHandler struct {
	carRepo         repository.CarRepository
	journeyRepo     repository.JourneyRepository
	assignerService *service.CarAssignerService
}

func NewDropoffHandler(carRepo repository.CarRepository, journeyRepo repository.JourneyRepository, assignerService *service.CarAssignerService) *DropoffHandler {
	return &DropoffHandler{
		carRepo:         carRepo,
		journeyRepo:     journeyRepo,
		assignerService: assignerService,
	}
}

func (h *DropoffHandler) Handle(dropoffRequest model.DropoffRequest) error {
	journey, err := h.journeyRepo.FindJourneyByID(dropoffRequest.ID)
	if err != nil {
		return fmt.Errorf("error finding journey: %v", err)
	}

	if journey.CarId != nil {
		car, err := h.carRepo.FindCarByID(*journey.CarId)
		if err != nil {
			return fmt.Errorf("error finding car for journey: %v", err)
		}
		h.assignerService.MoveCarToQueue(car, journey)
	}

	err = h.journeyRepo.DequeueJourney(journey.ID)
	if err != nil {
		return fmt.Errorf("error removing journey: %v", err)
	}

	return nil
}
