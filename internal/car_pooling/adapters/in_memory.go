package adapters

import (
	"car-pooling-service/internal/domain"
	"fmt"
	"sync"
	"time"
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

/////

type InMemoryJourneyRepository struct {
	journeys map[int]*dto.Journey
	mutex    sync.RWMutex
}

func NewInMemoryJourneyRepository() *InMemoryJourneyRepository {
	return &InMemoryJourneyRepository{
		journeys: make(map[int]*dto.Journey),
	}
}

func (r *InMemoryJourneyRepository) AddJourney(journey *dto.Journey) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.journeys[journey.ID]; exists {
		return fmt.Errorf("journey with ID %d already exists", journey.ID)
	}

	r.journeys[journey.ID] = journey

	return nil
}

func (r *InMemoryJourneyRepository) getAllJourneys() []*dto.Journey {
	var journeys []*dto.Journey
	for _, journey := range r.journeys {
		journeys = append(journeys, journey)
	}

	return journeys
}

func (r *InMemoryJourneyRepository) LogAllJourneys() {
	for _, journey := range r.journeys {
		fmt.Printf("Journey ID: %d, People: %d, Car Id: %d, WaitingSince: %s\n", journey.ID, journey.People, journey.CarId, journey.WaitingSince.Format(time.DateTime))
	}
	fmt.Printf("\n")
}

func (r *InMemoryJourneyRepository) FindJourneyByID(journeyID int) (*dto.Journey, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	journey, exists := r.journeys[journeyID]
	if !exists {
		return nil, fmt.Errorf("no journey found for group ID %d", journeyID)
	}

	return journey, nil
}

func (r *InMemoryCarRepository) FindCarByID(carID int) (*dto.Car, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	car, exists := r.cars[carID]
	if !exists {
		return nil, fmt.Errorf("no car found with ID %d", carID)
	}

	return car, nil
}
