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
	AddCar     *command.AddCarHandler
	AddJourney *command.AddJourneyHandler
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
			AddCar:     command.NewAddCarHandler(carRepoImpl),
			AddJourney: command.NewAddJourneyHandler(journeyRepoImpl),
		},
		Queries: Queries{
			Status:        query.NewStatusHandler(),
			LocateJourney: query.NewLocateJourneyHandler(journeyRepoImpl, carRepoImpl),
		},
	}
}
