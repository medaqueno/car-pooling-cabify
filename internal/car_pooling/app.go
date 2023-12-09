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
	Car *command.AddCarHandler
}

type Queries struct {
	Status *query.StatusHandler
}

func InitializeApp() *Application {
	// Prepare dependencies to be injected
	carRepo := adapters.NewInMemoryCarRepository()

	return &Application{
		Commands: Commands{
			Car: command.NewAddCarHandler(carRepo),
		},
		Queries: Queries{
			Status: query.NewStatusHandler(),
		},
	}
}
