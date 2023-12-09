package dto

import "time"

type AddCarRequest struct {
	ID    int `json:"id"`
	Seats int `json:"seats"`
}

type Car struct {
	ID             int
	Seats          int
	AvailableSeats int
}

func (c AddCarRequest) IsValid() bool {
	return c.Seats >= 4 && c.Seats <= 6
}

func NewCar(ID int, seats int, availableSeats int) *Car {
	return &Car{
		ID:             ID,
		Seats:          seats,
		AvailableSeats: availableSeats,
	}
}

type AddJourneyRequest struct {
	ID     int `json:"id"`
	People int `json:"people"`
}

type Journey struct {
	ID           int
	People       int
	CarId        *int
	WaitingSince time.Time
}

func (c AddJourneyRequest) IsValid() bool {
	return c.People >= 1 && c.People <= 6
}

func NewJourney(ID int, people int) *Journey {
	return &Journey{
		ID:           ID,
		People:       people,
		CarId:        nil,
		WaitingSince: time.Now(),
	}
}
