package repository

import "car-pooling-service/internal/domain"

type CarRepository interface {
	AddCar(car *dto.Car) error
	FindCarByID(carID int) (*dto.Car, error)
	GetAllCars() []*dto.Car
	LogAllCars()
}

type JourneyRepository interface {
	EnqueueJourney(car *dto.Journey) error
	FindJourneyByID(groupID int) (*dto.Journey, error)
	GetPendingJourneys() []*dto.Journey
	AssignCarToJourney(car *dto.Car, journey *dto.Journey) error
	RemoveJourney(carID int) error
	LogAllJourneys()
}
