package memory

import (
	"car-pooling-service/internal/domain/model"
	"container/heap"
	"errors"
	"fmt"
	"log"
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

type CarAssignerImpl struct {
	cars     map[int]*model.Car
	journeys JourneyQueue
	queues   [][]int // Each index represents the number of available seats
	mu       sync.Mutex
}

func NewCarAssignerRepository() *CarAssignerImpl {
	return &CarAssignerImpl{
		cars:     make(map[int]*model.Car),
		journeys: make(JourneyQueue, 0),
		queues:   make([][]int, 7), // Depends on Business Invariants -> Max Car Size + 1
	}
}

// Car
func (s *CarAssignerImpl) AddCar(car *model.Car) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.cars[car.ID]; exists {
		return fmt.Errorf("car with ID %d already exists", car.ID)
	}
	s.cars[car.ID] = car
	log.Printf("Add car: %d with %d\n", car.ID, car.Seats)

	err := s.addCarToQueue(car)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (s *CarAssignerImpl) addCarToQueue(car *model.Car) error {
	queueIndex := car.Seats

	s.queues[queueIndex] = append(s.queues[queueIndex], car.ID)
	car.InQueue = queueIndex

	// Update InQueue change in Car
	err := s.updateCar(car)
	if err != nil {
		return err
	}
	return nil
}

func (s *CarAssignerImpl) updateCar(car *model.Car) error {
	if _, exists := s.cars[car.ID]; !exists {
		return fmt.Errorf("car with ID %d not found", car.ID)
	}
	s.cars[car.ID] = car
	return nil
}

func (s *CarAssignerImpl) EnqueueJourney(journey *model.Journey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	heap.Push(&s.journeys, journey)
	log.Printf("Add Journey: %d with %d\n", journey.ID, journey.People)
	return nil
}

func (s *CarAssignerImpl) FindJourneyByID(journeyID int) (*model.Journey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, journey := range s.journeys {
		if journey.ID == journeyID {
			return journey, nil
		}
	}
	return nil, fmt.Errorf("journey not found") // Journey not found
}

func (s *CarAssignerImpl) FindCarByJourneyID(journeyID int) (*model.Car, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, journey := range s.journeys {
		if journey.ID == journeyID {
			if journey.CarId == nil {
				return nil, nil // No car assigned yet
			}
			return s.FindCarByID(*journey.CarId)
		}
	}

	return nil, fmt.Errorf("journey not found") // Journey not found
}

func (s *CarAssignerImpl) FindCarByID(carID int) (*model.Car, error) {
	if car, exists := s.cars[carID]; exists {
		return car, nil
	}

	return nil, fmt.Errorf("no car found with ID %d", carID)
}

func (s *CarAssignerImpl) DequeueJourney(journeyID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	/*if journey.CarId != nil {
		car, err := h.repo.FindCarByID(*journey.CarId)
		if err != nil {
			return fmt.Errorf("error finding car for journey: %v", err)
		}
		// h.repo.MoveCarToQueue(car, journey)

		
	}
	*/
	for i, journey := range s.journeys {
		if journey.ID == journeyID {
			s.journeys.Swap(i, s.journeys.Len()-1)
			s.journeys = s.journeys[:s.journeys.Len()-1]
			heap.Init(&s.journeys)
			return nil
		}
	}

	return errors.New("journey not found")
}

/*


func (s *CarAssignerImpl) AssignCarsToJourneys() {
	s.mu.Lock()
	defer s.mu.Unlock()

	waitingJourneys, err := s.journeyRepo.GetWaitingJourneys()
	if err != nil {
		log.Printf("Error fetching waiting journeys: %v", err)
		return
	}

	// PrintJourneys(waitingJourneys)

	for _, journey := range waitingJourneys {
		assigned := s.tryAssignCarToJourney(journey)
		if !assigned {
			// If we couldn't assign a car, the journey remains in the waiting list.
			// This is a good place to add logic for handling long-waiting journeys if we had Business Requirements
		}
	}

	// s.PrintAllQueues()
}

func (s *CarAssignerImpl) tryAssignCarToJourney(journey *model.Journey) bool {
	for seatsAvailable := journey.People; seatsAvailable < len(s.queues); seatsAvailable++ {
		if len(s.queues[seatsAvailable]) > 0 {
			// Found a car with enough seats
			carID := s.queues[seatsAvailable][0] // Get the first car in the queue

			car, err := s.carRepo.FindCarByID(carID)
			if err != nil {
				log.Printf("Error finding car: %v", err)
				return false
			}

			s.queues[seatsAvailable] = s.queues[seatsAvailable][1:] // Remove the car from the queue

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

			// log.Printf("Journey %d with %d people is assigned carID: %d that had %d seats, seatsNowAvailable: %d\n",
			//	journey.ID, journey.People, carID, car.Seats, seatsNowAvailable)

			return true
		}
	}
	return false
}

func (s *CarAssignerImpl) AddCarToQueue(car *model.Car) {
	s.mu.Lock()
	defer s.mu.Unlock()

	queueIndex := car.Seats

	s.queues[queueIndex] = append(s.queues[queueIndex], car.ID)
	car.InQueue = queueIndex

	// Update InQueue change in Car
	err := s.carRepo.UpdateCar(car)
	if err != nil {
		log.Printf("Error updating car: %v", err)
		return
	}
}

func (s *CarAssignerImpl) MoveCarToQueue(car *model.Car, journey *model.Journey) {
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

func removeCarFromQueue(queue []int, carID int) []int {
	for i, id := range queue {
		if id == carID {
			return append(queue[:i], queue[i+1:]...)
		}
	}
	return queue // In case car is not found, which should ideally not happen
}

func (s *CarAssignerImpl) PrintAllQueues() {

	for queueIndex, queue := range s.queues {
		fmt.Printf("Queue %d: , length: %d ", queueIndex, len(queue))
		fmt.Printf("| CarIDs:")
		for _, carID := range queue {
			fmt.Printf("%d", carID)
		}
		fmt.Println() // Newline for each queue
	}
	fmt.Println()
}

func PrintJourneys(journeys []*model.Journey) {
	for _, journey := range journeys {
		fmt.Printf("Journey ID: %d, People: %d, Waiting Since: %s\n",
			journey.ID, journey.People, journey.WaitingSince.Format(time.RFC3339))
	}
}
*/
