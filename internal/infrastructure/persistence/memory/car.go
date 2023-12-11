package memory

import (
	"car-pooling-service/internal/domain/model"
	"fmt"
	"sync"
)

type InMemoryCarRepository struct {
	cars  map[int]*model.Car
	mutex sync.RWMutex
}

func NewInMemoryCarRepository() *InMemoryCarRepository {
	return &InMemoryCarRepository{
		cars: make(map[int]*model.Car),
	}
}

func (r *InMemoryCarRepository) AddCar(car *model.Car) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.cars[car.ID]; exists {
		return fmt.Errorf("car with ID %d already exists", car.ID)
	}

	r.cars[car.ID] = car

	return nil
}

func (r *InMemoryCarRepository) GetAllCars() []*model.Car {
	var cars []*model.Car
	for _, car := range r.cars {
		cars = append(cars, car)
	}

	return cars
}

func (r *InMemoryCarRepository) FindCarByID(carID int) (*model.Car, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	car, exists := r.cars[carID]
	if !exists {
		return nil, fmt.Errorf("no car found with ID %d", carID)
	}

	return car, nil
}

func (r *InMemoryCarRepository) LogAllCars() {
	for _, car := range r.cars {
		fmt.Printf("Car ID: %d, Seats: %d, Available Seats: %d\n", car.ID, car.Seats, car.AvailableSeats)
	}
	fmt.Printf("\n")
}
