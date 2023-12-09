package app

import (
	"car-pooling-service/internal/car_pooling/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	Status *query.StatusHandler
}

func InitializeApp() *Application {
	return &Application{
		Commands: Commands{},
		Queries: Queries{
			Status: query.NewStatusHandler(),
		},
	}
}
