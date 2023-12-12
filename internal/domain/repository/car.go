package repository

import "car-pooling-service/internal/domain/model"

type CarRepository interface {
	AddCar(car *model.Car) error
	FindCarByID(carID int) (*model.Car, error)
	UpdateCar(car *model.Car) error
}
