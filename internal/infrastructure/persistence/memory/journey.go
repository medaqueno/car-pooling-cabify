package memory

import (
	"car-pooling-service/internal/domain/model"
	"errors"
	"sync"
)

type JourneyRepositoryImpl struct {
	journeys map[int]*model.Journey
	mu       sync.Mutex
}

func NewJourneyRepository() *JourneyRepositoryImpl {
	return &JourneyRepositoryImpl{
		journeys: make(map[int]*model.Journey),
	}
}

func (r *JourneyRepositoryImpl) EnqueueJourney(journey *model.Journey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.journeys[journey.ID]; exists {
		return errors.New("journey already exists")
	}
	r.journeys[journey.ID] = journey
	return nil
}

func (r *JourneyRepositoryImpl) DequeueJourney(journeyID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.journeys[journeyID]; !exists {
		return errors.New("journey not found")
	}
	delete(r.journeys, journeyID)
	return nil
}

func (r *JourneyRepositoryImpl) UpdateJourney(journey *model.Journey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.journeys[journey.ID] = journey
	return nil
}

func (r *JourneyRepositoryImpl) FindJourneyByID(journeyID int) (*model.Journey, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if journey, exists := r.journeys[journeyID]; exists {
		return journey, nil
	}
	return nil, errors.New("journey not found")
}

func (r *JourneyRepositoryImpl) GetWaitingJourneys() ([]*model.Journey, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var waitingJourneys []*model.Journey
	for _, journey := range r.journeys {
		if journey.CarId == nil {
			waitingJourneys = append(waitingJourneys, journey)
		}
	}

	return waitingJourneys, nil
}
