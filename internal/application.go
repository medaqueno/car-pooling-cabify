package internal

import (
	"car-pooling-service/internal/application/command"
	"car-pooling-service/internal/application/query"
	"car-pooling-service/internal/application/service"
	"car-pooling-service/internal/infrastructure/persistence/memory"
)

type Application struct {
	Commands Commands
	Queries  Queries
	Services Services
}

type Commands struct {
	AddCar         *command.AddCarHandler
	EnqueueJourney *command.EnqueueJourneyHandler
	Dropoff        *command.DropoffHandler
}

type Queries struct {
	Status        *query.StatusHandler
	LocateJourney *query.LocateCarByJourneyHandler
}

type Services struct {
	AssignCarsToJourneys *service.AssignCarsToJourneysHandler
}

func InitializeApp() *Application {
	// Prepare dependencies to be injected
	carRepoImpl := memory.NewInMemoryCarRepository()
	journeyRepoImpl := memory.NewInMemoryJourneyRepository()

	return &Application{
		Commands: Commands{
			AddCar:         command.NewAddCarHandler(carRepoImpl),
			EnqueueJourney: command.NewEnqueueJourneyHandler(journeyRepoImpl),

			Dropoff: command.NewDropoffHandler(carRepoImpl, journeyRepoImpl),
		},
		Queries: Queries{
			Status:        query.NewStatusHandler(),
			LocateJourney: query.NewLocateJourneyHandler(carRepoImpl, journeyRepoImpl),
		},
		Services: Services{
			AssignCarsToJourneys: service.NewAssignCarsToJourneysHandler(carRepoImpl, journeyRepoImpl),
		},
	}
}
