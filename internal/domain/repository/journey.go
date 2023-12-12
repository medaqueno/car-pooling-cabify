package repository

import "car-pooling-service/internal/domain/model"

type JourneyRepository interface {
	EnqueueJourney(journey *model.Journey) error
	DequeueJourney(journeyID int) error
	UpdateJourney(journey *model.Journey) error
	FindJourneyByID(journeyID int) (*model.Journey, error)
	GetWaitingJourneys() ([]*model.Journey, error)
}
