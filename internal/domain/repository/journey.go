package repository

import "car-pooling-service/internal/domain/model"

type JourneyRepository interface {
	EnqueueJourney(car *model.Journey) error
	FindJourneyByID(groupID int) (*model.Journey, error)
	GetPendingJourneys() []*model.Journey
	AssignCarToJourney(car *model.Car, journey *model.Journey) error
	RemoveJourney(carID int) error
	LogAllJourneys()
}
