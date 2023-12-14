package repository

import "car-pooling-service/internal/domain/model"

type CarAssignerRepository interface {
	AssignCarsToJourneys()
	AddCarToQueue(car *model.Car)
	MoveCarToQueue(car *model.Car, journey *model.Journey)
}
