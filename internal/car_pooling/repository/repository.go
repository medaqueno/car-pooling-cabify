package repository

import "car-pooling-service/internal/domain"

type CarRepository interface {
	AddCar(car *dto.Car) error
	FindCarByID(carID int) (*dto.Car, error)
	LogAllCars()
}

type JourneyRepository interface {
	AddJourney(car *dto.Journey) error
	FindJourneyByID(groupID int) (*dto.Journey, error)
	LogAllJourneys()
}
