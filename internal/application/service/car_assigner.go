// service/car_assigner.go

package service

import (
	"car-pooling-service/internal/domain/model"
	"car-pooling-service/internal/domain/repository"
	"fmt"
	"log"
	"sync"
	"time"
)

type CarAssignerService struct {
	carRepo     repository.CarRepository
	journeyRepo repository.JourneyRepository
	queues      [][]int // Each index represents the number of available seats
	mu          sync.Mutex
}

func NewCarAssignerService(carRepo repository.CarRepository, journeyRepo repository.JourneyRepository, queueCount int) *CarAssignerService {
	queues := make([][]int, queueCount)
	for i := range queues {
		queues[i] = make([]int, 0)
	}

	return &CarAssignerService{
		carRepo:     carRepo,
		journeyRepo: journeyRepo,
		queues:      queues,
	}
}

func (s *CarAssignerService) RunAssignmentProcess() {
	go func() {
		for {
			s.AssignCarsToJourneys()
			time.Sleep(5 * time.Second) // Run every 5 seconds // TODO: Move to Event Listener
		}
	}()
}

func (s *CarAssignerService) AssignCarsToJourneys() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Fetch waiting journeys
	waitingJourneys, err := s.journeyRepo.GetWaitingJourneys()
	if err != nil {
		log.Printf("Error fetching waiting journeys: %v", err)
		return
	}

	// s.PrintAllQueues()

	for _, journey := range waitingJourneys {
		assigned := s.tryAssignCarToJourney(journey)
		if !assigned {
			// If we couldn't assign a car, the journey remains in the waiting list.
			// Consider adding logic for handling long-waiting journeys.
		}
	}
}

func (s *CarAssignerService) tryAssignCarToJourney(journey *model.Journey) bool {

	for seatsAvailable := journey.People; seatsAvailable < len(s.queues); seatsAvailable++ {
		if len(s.queues[seatsAvailable]) > 0 {

			// Found a car with enough seats
			carID := s.queues[seatsAvailable][0]                    // Get the first car in the queue
			s.queues[seatsAvailable] = s.queues[seatsAvailable][1:] // Remove the car from the queue

			car, err := s.carRepo.FindCarByID(carID)
			if err != nil {
				log.Printf("Error finding car: %v", err)
				continue
			}

			// Update journey with the car's ID
			journey.CarId = &carID
			err = s.journeyRepo.UpdateJourney(journey)
			if err != nil {
				log.Printf("Error updating journey: %v", err)
				return false
			}

			// Update car's availability and InQueue attribute
			seatsNowAvailable := seatsAvailable - journey.People
			s.queues[seatsNowAvailable] = append(s.queues[seatsNowAvailable], carID)
			car.InQueue = seatsNowAvailable

			err = s.carRepo.UpdateCar(car)
			if err != nil {
				log.Printf("Error updating car: %v", err)
				return false
			}

			return true
		}
	}
	return false
}

func (s *CarAssignerService) AddCarToQueue(car *model.Car) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queueIndex := car.Seats // Assuming this represents the number of AVAILABLE seats
	s.queues[queueIndex] = append(s.queues[queueIndex], car.ID)
	car.InQueue = queueIndex

	// Update InQueue change in Car
	err := s.carRepo.UpdateCar(car)
	if err != nil {
		log.Printf("Error updating car: %v", err)
		return
	}
}

func (s *CarAssignerService) MoveCarToQueue(car *model.Car, journey *model.Journey) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove car from its current queue
	currentQueue := car.InQueue
	s.queues[currentQueue] = removeCarFromQueue(s.queues[currentQueue], car.ID)

	// Calculate new Queue = Current Available seats + total of People dropping off
	newQueueIndex := car.InQueue + journey.People

	// Add car to the new queue
	s.queues[newQueueIndex] = append(s.queues[newQueueIndex], car.ID)
	car.InQueue = newQueueIndex

	err := s.carRepo.UpdateCar(car)
	if err != nil {
		log.Printf("Error updating car: %v", err)
		return
	}
}

// Helper function to remove a car from a queue
func removeCarFromQueue(queue []int, carID int) []int {
	for i, id := range queue {
		if id == carID {
			return append(queue[:i], queue[i+1:]...)
		}
	}
	return queue // In case car is not found, which should ideally not happen
}

func (s *CarAssignerService) PrintAllQueues() {

	for queueIndex, queue := range s.queues {
		fmt.Printf("Queue %d: , length: %d ", queueIndex, len(queue))
		fmt.Printf("| CarIDs:")
		for _, carID := range queue {
			fmt.Printf("%d ", carID)
		}
		fmt.Println() // Newline for each queue
	}
	fmt.Println()
}
