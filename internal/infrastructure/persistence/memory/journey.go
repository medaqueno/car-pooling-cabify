package memory

import (
	"car-pooling-service/internal/domain/model"
	"container/heap"
	"errors"
	"sync"
)

type JourneyQueue []*model.Journey

func (jq JourneyQueue) Len() int { return len(jq) }
func (jq JourneyQueue) Less(i, j int) bool {
	// First journeys to enter in queue first. FIFO-like
	return jq[i].WaitingSince.Before(jq[j].WaitingSince)
}
func (jq JourneyQueue) Swap(i, j int) { jq[i], jq[j] = jq[j], jq[i] }

func (jq *JourneyQueue) Push(x interface{}) {
	item := x.(*model.Journey)
	*jq = append(*jq, item)
}

func (jq *JourneyQueue) Pop() interface{} {
	old := *jq
	n := len(old)
	item := old[n-1]
	*jq = old[0 : n-1]
	return item
}

type JourneyRepositoryImpl struct {
	journeys JourneyQueue
	mu       sync.Mutex
}

func NewJourneyRepository() *JourneyRepositoryImpl {
	return &JourneyRepositoryImpl{
		journeys: make(JourneyQueue, 0),
	}
}

func (r *JourneyRepositoryImpl) EnqueueJourney(journey *model.Journey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	heap.Push(&r.journeys, journey)
	return nil
}

func (r *JourneyRepositoryImpl) DequeueJourney(journeyID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, journey := range r.journeys {
		if journey.ID == journeyID {
			r.journeys.Swap(i, r.journeys.Len()-1)
			r.journeys = r.journeys[:r.journeys.Len()-1]
			heap.Init(&r.journeys)
			return nil
		}
	}

	return errors.New("journey not found")
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

func (r *JourneyRepositoryImpl) UpdateJourney(updatedJourney *model.Journey) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, journey := range r.journeys {
		if journey.ID == updatedJourney.ID {
			r.journeys[i] = updatedJourney

			// No need to re-init Heap because we modified a non-ordered attribute
			return nil
		}
	}

	return errors.New("journey not found")
}

func (r *JourneyRepositoryImpl) FindJourneyByID(journeyID int) (*model.Journey, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, journey := range r.journeys {
		if journey.ID == journeyID {
			return journey, nil
		}
	}

	return nil, errors.New("journey not found")
}
