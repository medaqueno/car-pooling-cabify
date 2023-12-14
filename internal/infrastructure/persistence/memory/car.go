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

func (repo *CarRepositoryImpl) AddCar(car *model.Car) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.cars[car.ID]; exists {
		// return errors.New("car already exists")
		return fmt.Errorf("car with ID %d already exists", car.ID)
	}
	repo.cars[car.ID] = car
	return nil
}

func (repo *CarRepositoryImpl) FindCarByID(carID int) (*model.Car, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if car, exists := repo.cars[carID]; exists {
		return car, nil
	}
	// return nil, errors.New("car not found")
	return nil, fmt.Errorf("no car found with ID %d", carID)
}

func (repo *CarRepositoryImpl) UpdateCar(car *model.Car) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.cars[car.ID]; !exists {
		return fmt.Errorf("car with ID %d not found", car.ID)
	}
	repo.cars[car.ID] = car
	return nil
}
