package command

import (
	"car-pooling-service/internal/car_pooling/repository"
	"car-pooling-service/internal/domain"
)

type AssignCarsToJourneysHandler struct {
	carRepo     repository.CarRepository
	journeyRepo repository.JourneyRepository
}

func NewAssignCarsToJourneysHandler(carRepo repository.CarRepository, journeyRepo repository.JourneyRepository) *AssignCarsToJourneysHandler {
	return &AssignCarsToJourneysHandler{
		carRepo:     carRepo,
		journeyRepo: journeyRepo,
	}
}

func (h *AssignCarsToJourneysHandler) Handle() error {
	pendingJourneys := h.journeyRepo.GetPendingJourneys()
	availableCars := h.carRepo.GetAllCars()

	for _, journey := range pendingJourneys {
		for _, car := range availableCars {
			if h.isSuitableCar(journey, car) {
				err := h.journeyRepo.AssignCarToJourney(car, journey)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return nil
}

// isSuitableCar Check if the car has enough seats available for the group
// and if the car is not already occupied
func (h *AssignCarsToJourneysHandler) isSuitableCar(journey *dto.Journey, car *dto.Car) bool {
	return car.AvailableSeats >= journey.People && car.AvailableSeats == car.Seats
}
