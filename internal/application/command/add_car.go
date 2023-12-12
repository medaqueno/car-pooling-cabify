package command

import (
	"car-pooling-service/internal/application/service"
	"car-pooling-service/internal/domain/model"
	"car-pooling-service/internal/domain/repository"
	"fmt"
)

type AddCarHandler struct {
	repo            repository.CarRepository
	assignerService *service.CarAssignerService
}

func NewAddCarHandler(repo repository.CarRepository, assignerService *service.CarAssignerService) *AddCarHandler {
	return &AddCarHandler{
		repo:            repo,
		assignerService: assignerService,
	}
}

func (h *AddCarHandler) Handle(addCarsRequest []model.AddCarRequest) error {
	for _, carRequest := range addCarsRequest {
		car := model.NewCar(carRequest.ID, carRequest.Seats)

		// Add car to the repository
		err := h.repo.AddCar(car)
		if err != nil {
			fmt.Printf("Error adding car: %v\n", err)
			continue
		}

		// Add car to the corresponding AvailabilityQueue
		h.assignerService.AddCarToQueue(car)
	}
	
	return nil
}
