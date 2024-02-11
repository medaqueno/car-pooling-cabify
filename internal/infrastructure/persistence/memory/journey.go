package memory

import (
	"car-pooling-service/internal/domain/model"
	"container/heap"
	"errors"
	"fmt"
	"sync"
)

type JourneyQueue []*model.Journey

func (jq *JourneyQueue) Len() int { return len(*jq) }
func (jq *JourneyQueue) Less(i, j int) bool {
	// First journeys to enter in queue first. FIFO-like
	return (*jq)[i].WaitingSince.Before((*jq)[j].WaitingSince)
}
func (jq *JourneyQueue) Swap(i, j int) { (*jq)[i], (*jq)[j] = (*jq)[j], (*jq)[i] }

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

func (repo *JourneyRepositoryImpl) EnqueueJourney(journey *model.Journey) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	heap.Push(&repo.journeys, journey)
	return nil
}

func (repo *JourneyRepositoryImpl) DequeueJourney(journeyID int) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i, journey := range repo.journeys {
		if journey.ID == journeyID {
			repo.journeys.Swap(i, repo.journeys.Len()-1)
			repo.journeys = repo.journeys[:repo.journeys.Len()-1]
			heap.Init(&repo.journeys)
			return nil
		}
	}

	return errors.New("journey not found")
}

func (repo *JourneyRepositoryImpl) GetWaitingJourneys() ([]*model.Journey, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var waitingJourneys []*model.Journey
	for _, journey := range repo.journeys {
		if journey.CarId == nil {
			waitingJourneys = append(waitingJourneys, journey)
		}
	}
	return waitingJourneys, nil
}

func (repo *JourneyRepositoryImpl) UpdateJourney(updatedJourney *model.Journey) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i, journey := range repo.journeys {
		if journey.ID == updatedJourney.ID {
			repo.journeys[i] = updatedJourney

			// No need to re-init Heap because we modified a non-ordered attribute
			return nil
		}
	}

	return errors.New("journey not found")
}

func (repo *JourneyRepositoryImpl) FindJourneyByID(journeyID int) (*model.Journey, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, journey := range repo.journeys {
		if journey.ID == journeyID {
			return journey, nil
		}
	}

	return nil, fmt.Errorf("journey not found")
}
