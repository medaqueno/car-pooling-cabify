package repository

import "car-pooling-service/internal/domain/model"

type CarAssignerRepository interface {
	AddCar(car *model.Car) error
	EnqueueJourney(journey *model.Journey) error
	FindCarByJourneyID(journeyID int) (*model.Car, error)
	DequeueJourney(journeyID int) error
	FindJourneyByID(journeyID int) (*model.Journey, error)
	FindCarByID(carID int) (*model.Car, error)
	//AddCarToQueue(car *model.Car)
	/*
		AssignCarsToJourneys()
		AddCarToQueue(car *model.Car)
		MoveCarToQueue(car *model.Car, journey *model.Journey)
	*/
}
