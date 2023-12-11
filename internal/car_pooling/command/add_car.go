package command

import (
	"car-pooling-service/internal/car_pooling/repository"
	"car-pooling-service/internal/domain"
	"fmt"
)

type AddCarHandler struct {
	repo repository.CarRepository
}

func NewAddCarHandler(repo repository.CarRepository) *AddCarHandler {
	return &AddCarHandler{
		repo: repo,
	}
}

func (h *AddCarHandler) Handle(addCarsRequest []dto.AddCarRequest) error {
	for _, carRequest := range addCarsRequest {
		car := dto.NewCar(carRequest.ID, carRequest.Seats, carRequest.Seats)
		err := h.repo.AddCar(car)
		if err != nil {
			fmt.Printf("Error adding car: %v\n", err)
			continue
		}
	}

	// Debug
	// h.repo.LogAllCars()

	return nil
}
