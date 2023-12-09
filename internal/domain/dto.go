package dto

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

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

type LocateJourneyRequest struct {
	ID int
}

func (r *LocateJourneyRequest) Validate(req *http.Request) error {
	err := req.ParseForm()

	if err != nil {
		return fmt.Errorf("bad request")
	}

	idStr := req.Form.Get("ID")
	if idStr == "" {
		return fmt.Errorf("ID is required")
	}

	r.ID, err = strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("ID must be an integer")
	}

	if r.ID <= 0 {
		return fmt.Errorf("ID must be a positive integer")
	}

	return nil
}
