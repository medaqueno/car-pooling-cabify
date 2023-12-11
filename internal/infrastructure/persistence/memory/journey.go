package memory

import (
	"car-pooling-service/internal/domain/model"
	"fmt"
	"sync"
	"time"
)

type InMemoryJourneyRepository struct {
	journeys []*model.Journey
	mutex    sync.RWMutex
}

func NewInMemoryJourneyRepository() *InMemoryJourneyRepository {
	return &InMemoryJourneyRepository{
		journeys: []*model.Journey{},
	}
}

func (r *InMemoryJourneyRepository) EnqueueJourney(journey *model.Journey) error {
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

func (r *InMemoryJourneyRepository) GetPendingJourneys() []*model.Journey {
	var journeys []*model.Journey
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
func (r *InMemoryJourneyRepository) FindJourneyByID(journeyID int) (*model.Journey, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, journey := range r.journeys {
		if journey.ID == journeyID {
			return journey, nil
		}
	}

	return nil, fmt.Errorf("no journey found for group ID %d", journeyID)
}

func (r *InMemoryJourneyRepository) AssignCarToJourney(car *model.Car, journey *model.Journey) error {
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
