package service

import (
	"car-pooling-service/internal/domain/repository"
)

type CarAssignerHandler struct {
	carAssignerRepo repository.CarAssignerRepository
}

func NewCarAssignerHandler(carAssignerRepo repository.CarAssignerRepository) *CarAssignerHandler {
	return &CarAssignerHandler{
		carAssignerRepo: carAssignerRepo,
	}
}

func (s *CarAssignerHandler) Handle() {
	s.carAssignerRepo.AssignCarsToJourneys()
}
