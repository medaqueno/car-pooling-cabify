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

func (r *InMemoryCarRepository) GetAllCars() []*dto.Car {
	var cars []*dto.Car
	for _, car := range r.cars {
		cars = append(cars, car)
	}

	return cars
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

func (r *InMemoryCarRepository) LogAllCars() {
	for _, car := range r.cars {
		fmt.Printf("Car ID: %d, Seats: %d, Available Seats: %d\n", car.ID, car.Seats, car.AvailableSeats)
	}
	fmt.Printf("\n")
}

/////

type InMemoryJourneyRepository struct {
	journeys []*dto.Journey
	mutex    sync.RWMutex
}

func NewInMemoryJourneyRepository() *InMemoryJourneyRepository {
	return &InMemoryJourneyRepository{
		journeys: []*dto.Journey{},
	}
}

func (r *InMemoryJourneyRepository) EnqueueJourney(journey *dto.Journey) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, j := range r.journeys {
		if j.ID == journey.ID {
			return fmt.Errorf("journey with ID %d already exists", journey.ID)
		}
	}

	r.journeys = append(r.journeys, journey)

	return nil
}

func (r *InMemoryJourneyRepository) GetPendingJourneys() []*dto.Journey {
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

	for _, journey := range r.journeys {
		if journey.ID == journeyID {
			return journey, nil
		}
	}

	return nil, fmt.Errorf("no journey found for group ID %d", journeyID)
}

func (r *InMemoryJourneyRepository) AssignCarToJourney(car *dto.Car, journey *dto.Journey) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, existingJourney := range r.journeys {
		if existingJourney.ID == journey.ID {
			if existingJourney.CarId != nil {
				return nil
			}

			existingJourney.CarId = &car.ID
			car.AvailableSeats -= journey.People
			if car.AvailableSeats < 0 {
				car.AvailableSeats = 0
			}
			fmt.Printf("\nAssigned Journey to Car %d\n", existingJourney.ID)
			return nil
		}
	}

	return fmt.Errorf("journey with ID %d not found", journey.ID)
}

func (r *InMemoryJourneyRepository) RemoveJourney(journeyID int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i, journey := range r.journeys {
		if journey.ID == journeyID {
			r.journeys = append(r.journeys[:i], r.journeys[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("journey with ID %d not found", journeyID)
}
