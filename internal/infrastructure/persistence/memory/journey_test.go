package memory

import (
	"car-pooling-service/internal/domain/model"
	"testing"
)

func TestEnqueueDequeueJourney(t *testing.T) {
	repo := NewJourneyRepository()
	journey := model.NewJourney(1, 4)

	if err := repo.EnqueueJourney(journey); err != nil {
		t.Errorf("EnqueueJourney failed: %v", err)
	}

	// Test enqueuing a duplicate journey
	if err := repo.EnqueueJourney(journey); err == nil {
		t.Error("Expected an error when enqueuing a duplicate journey, but got none")
	}

	if err := repo.DequeueJourney(1); err != nil {
		t.Errorf("DequeueJourney failed: %v", err)
	}

	// Test dequeuing a non-existing journey
	if err := repo.DequeueJourney(1); err == nil {
		t.Error("Expected an error when dequeuing a non-existing journey, but got none")
	}
}

func TestFindJourneyByID(t *testing.T) {
	repo := NewJourneyRepository()
	journey := model.NewJourney(1, 4)
	_ = repo.EnqueueJourney(journey)

	if _, err := repo.FindJourneyByID(1); err != nil {
		t.Errorf("FindJourneyByID failed to find journey: %v", err)
	}

	// Test finding a non-existing journey
	if _, err := repo.FindJourneyByID(2); err == nil {
		t.Error("Expected an error when finding a non-existing journey, but got none")
	}
}

func TestUpdateJourney(t *testing.T) {
	repo := NewJourneyRepository()
	journey := model.NewJourney(1, 4)
	_ = repo.EnqueueJourney(journey)

	carID := 123
	journey.CarId = &carID

	if err := repo.UpdateJourney(journey); err != nil {
		t.Errorf("UpdateJourney failed: %v", err)
	}

	updatedJourney, _ := repo.FindJourneyByID(1)
	if updatedJourney.CarId == nil || *updatedJourney.CarId != carID {
		t.Errorf("UpdateJourney did not correctly update the CarId")
	}

	// Test updating a non-existing journey
	nonExistingJourney := model.NewJourney(2, 4)
	if err := repo.UpdateJourney(nonExistingJourney); err == nil {
		t.Error("Expected an error when updating a non-existing journey, but got none")
	}
}

func TestGetWaitingJourneys(t *testing.T) {
	repo := NewJourneyRepository()
	journey1 := model.NewJourney(1, 4)
	journey2 := model.NewJourney(2, 3)

	_ = repo.EnqueueJourney(journey1)
	_ = repo.EnqueueJourney(journey2)

	waitingJourneys, err := repo.GetWaitingJourneys()
	if err != nil {
		t.Errorf("GetWaitingJourneys failed: %v", err)
	}
	if len(waitingJourneys) != 2 {
		t.Errorf("Expected 2 waiting journeys, got %d", len(waitingJourneys))
	}

	// Dequeue one journey and test if everything was updated
	_ = repo.DequeueJourney(1)
	waitingJourneys, _ = repo.GetWaitingJourneys()
	if len(waitingJourneys) != 1 {
		t.Errorf("Expected 1 waiting journey after dequeue, got %d", len(waitingJourneys))
	}

	// Check the remaining journey's ID
	if waitingJourneys[0].ID != 2 {
		t.Errorf("Expected journey ID 2 to be waiting, got ID %d", waitingJourneys[0].ID)
	}
}
