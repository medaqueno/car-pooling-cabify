package memory

import (
	"car-pooling-service/internal/domain/model"
	"testing"
	"time"
)

func TestCarAssignerAssignCarsToJourneys(t *testing.T) {
	carRepo := NewCarRepository()
	journeyRepo := NewJourneyRepository()

	// Add Cars
	car1 := &model.Car{ID: 1, Seats: 4}
	_ = carRepo.AddCar(car1)

	// add journeys
	journey1 := &model.Journey{ID: 1, People: 2, WaitingSince: time.Now()}
	_ = journeyRepo.EnqueueJourney(journey1)

	carAssigner := NewCarAssignerRepository(carRepo, journeyRepo)

	// Add car
	carAssigner.AddCarToQueue(car1)

	carAssigner.AssignCarsToJourneys()

	updatedJourney, err := journeyRepo.FindJourneyByID(1)
	if err != nil {
		t.Fatalf("Failed to find journey: %v", err)
	}
	if updatedJourney.CarId == nil || *updatedJourney.CarId != 1 {
		t.Errorf("Expected journey to be assigned to car 1, got car %v", updatedJourney.CarId)
	}
}

func TestCarAssignerReassignmentAfterDropoff(t *testing.T) {
	carRepo := NewCarRepository()
	journeyRepo := NewJourneyRepository()
	carAssigner := NewCarAssignerRepository(carRepo, journeyRepo)

	// Add car and Journey
	car := &model.Car{ID: 1, Seats: 4}
	_ = carRepo.AddCar(car)
	journey := &model.Journey{ID: 1, People: 3, WaitingSince: time.Now()}
	_ = journeyRepo.EnqueueJourney(journey)

	// Add Car to queue
	carAssigner.AddCarToQueue(car)
	carAssigner.AssignCarsToJourneys()

	// Simulate dropoff
	carAssigner.MoveCarToQueue(car, journey)

	// Add new JOurney
	newJourney := &model.Journey{ID: 2, People: 2, WaitingSince: time.Now()}
	_ = journeyRepo.EnqueueJourney(newJourney)

	carAssigner.AssignCarsToJourneys()

	// Verify
	updatedJourney, err := journeyRepo.FindJourneyByID(2)
	if err != nil {
		t.Fatalf("Failed to find new journey: %v", err)
	}
	if updatedJourney.CarId == nil || *updatedJourney.CarId != 1 {
		t.Errorf("Expected new journey to be assigned to car 1, got car %v", updatedJourney.CarId)
	}
}

func TestCarAssignerPriorityAssignment(t *testing.T) {
	carRepo := NewCarRepository()
	journeyRepo := NewJourneyRepository()
	carAssigner := NewCarAssignerRepository(carRepo, journeyRepo)

	// Add a car
	car := &model.Car{ID: 1, Seats: 4}
	_ = carRepo.AddCar(car)
	carAssigner.AddCarToQueue(car)

	// ADd Journeys in different moments
	earlyJourney := &model.Journey{ID: 1, People: 3, WaitingSince: time.Now().Add(-10 * time.Minute)}
	_ = journeyRepo.EnqueueJourney(earlyJourney)

	time.Sleep(1 * time.Second) // Ensure new "later" journey

	lateJourney := &model.Journey{ID: 2, People: 2, WaitingSince: time.Now()}
	_ = journeyRepo.EnqueueJourney(lateJourney)

	carAssigner.AssignCarsToJourneys()

	// Verify that oldest journey was assigned first
	updatedEarlyJourney, err := journeyRepo.FindJourneyByID(1)
	if err != nil {
		t.Fatalf("Failed to find early journey: %v", err)
	}
	if updatedEarlyJourney.CarId == nil || *updatedEarlyJourney.CarId != 1 {
		t.Errorf("Expected early journey to be assigned to car 1, got car %v", updatedEarlyJourney.CarId)
	}

	// Verify that newest journey was not assigned yet
	updatedLateJourney, err := journeyRepo.FindJourneyByID(2)
	if err != nil {
		t.Fatalf("Failed to find late journey: %v", err)
	}
	if updatedLateJourney.CarId != nil {
		t.Errorf("Expected late journey to not be assigned yet, but was assigned to car %d", *updatedLateJourney.CarId)
	}
}
