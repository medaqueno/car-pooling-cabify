package command

import (
	"car-pooling-service/internal/domain/model"
	"car-pooling-service/internal/domain/repository"
	"car-pooling-service/internal/port/http/dto"
	"fmt"
)

type AddCarHandler struct {
	repo                  repository.CarRepository
	carAssignerRepository repository.CarAssignerRepository
}

func NewAddCarHandler(repo repository.CarRepository, carAssignerRepository repository.CarAssignerRepository) *AddCarHandler {
	return &AddCarHandler{
		repo:                  repo,
		carAssignerRepository: carAssignerRepository,
	}
}

func (h *AddCarHandler) Handle(addCarsRequest []dto.AddCarRequest) error {
	for _, carRequest := range addCarsRequest {
		car := model.NewCar(carRequest.ID, carRequest.Seats)

		// Add car to the repository
		err := h.repo.AddCar(car)
		if err != nil {
			fmt.Printf("Error adding car: %v\n", err)
			continue
		}

		// Add car to the corresponding AvailabilityQueue
		h.carAssignerRepository.AddCarToQueue(car)
	}

	return nil
}
