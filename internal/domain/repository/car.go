package repository

import "car-pooling-service/internal/domain/model"

type CarRepository interface {
	AddCar(car *model.Car) error
	FindCarByID(carID int) (*model.Car, error)
	GetAllCars() []*model.Car
	LogAllCars()
}
