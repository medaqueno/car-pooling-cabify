package memory

import (
	"car-pooling-service/internal/domain/model"
	"fmt"
	"sync"
)

type CarRepositoryImpl struct {
	cars map[int]*model.Car
	mu   sync.Mutex
}

func NewCarRepository() *CarRepositoryImpl {
	return &CarRepositoryImpl{
		cars: make(map[int]*model.Car),
	}
}

func (r *CarRepositoryImpl) AddCar(car *model.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.cars[car.ID]; exists {
		// return errors.New("car already exists")
		return fmt.Errorf("car with ID %d already exists", car.ID)
	}
	r.cars[car.ID] = car
	return nil
}

func (r *CarRepositoryImpl) FindCarByID(carID int) (*model.Car, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if car, exists := r.cars[carID]; exists {
		return car, nil
	}
	// return nil, errors.New("car not found")
	return nil, fmt.Errorf("no car found with ID %d", carID)
}

func (r *CarRepositoryImpl) UpdateCar(car *model.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cars[car.ID]; !exists {
		return fmt.Errorf("car with ID %d not found", car.ID)
	}
	r.cars[car.ID] = car
	return nil
}
