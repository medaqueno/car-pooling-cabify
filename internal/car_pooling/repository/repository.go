package repository

import "car-pooling-service/internal/domain"

type CarRepository interface {
	AddCar(car *dto.Car) error
	LogAllCars()
}
