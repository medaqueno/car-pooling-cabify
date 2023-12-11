package app

import (
	"car-pooling-service/internal/car_pooling/adapters"
	"car-pooling-service/internal/car_pooling/command"
	"car-pooling-service/internal/car_pooling/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	AddCar               *command.AddCarHandler
	EnqueueJourney       *command.EnqueueJourneyHandler
	AssignCarsToJourneys *command.AssignCarsToJourneysHandler
	Dropoff              *command.DropoffHandler
}

type Queries struct {
	Status        *query.StatusHandler
	LocateJourney *query.LocateCarByJourneyHandler
}

func InitializeApp() *Application {
	// Prepare dependencies to be injected
	carRepoImpl := adapters.NewInMemoryCarRepository()
	journeyRepoImpl := adapters.NewInMemoryJourneyRepository()

	return &Application{
		Commands: Commands{
			AddCar:               command.NewAddCarHandler(carRepoImpl),
			EnqueueJourney:       command.NewEnqueueJourneyHandler(journeyRepoImpl),
			AssignCarsToJourneys: command.NewAssignCarsToJourneysHandler(carRepoImpl, journeyRepoImpl),
			Dropoff:              command.NewDropoffHandler(carRepoImpl, journeyRepoImpl),
		},
		Queries: Queries{
			Status:        query.NewStatusHandler(),
			LocateJourney: query.NewLocateJourneyHandler(carRepoImpl, journeyRepoImpl),
		},
	}
}
