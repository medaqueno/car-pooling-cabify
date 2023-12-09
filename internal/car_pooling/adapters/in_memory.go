package adapters

import (
	"car-pooling-service/internal/domain"
	"fmt"
	"sync"
)

type InMemoryCarRepository struct {
	cars  map[int]*dto.Car
	mutex sync.RWMutex
}

func NewInMemoryCarRepository() *InMemoryCarRepository {
	return &InMemoryCarRepository{
		cars: make(map[int]*dto.Car),
	}
}

func (r *InMemoryCarRepository) AddCar(car *dto.Car) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.cars[car.ID]; exists {
		return fmt.Errorf("car with ID %d already exists", car.ID)
	}

	r.cars[car.ID] = car

	return nil
}

func (r *InMemoryCarRepository) getAllCars() []*dto.Car {
	var cars []*dto.Car
	for _, car := range r.cars {
		cars = append(cars, car)
	}

	return cars
}

func (r *InMemoryCarRepository) LogAllCars() {
	for _, car := range r.cars {
		fmt.Printf("Car ID: %d, Seats: %d, Available Seats: %d\n", car.ID, car.Seats, car.AvailableSeats)
	}
	fmt.Printf("\n")
}
