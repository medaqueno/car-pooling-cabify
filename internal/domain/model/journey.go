package model

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type EnqueueJourneyRequest struct {
	ID     int `json:"id"`
	People int `json:"people"`
}

type Journey struct {
	ID           int
	People       int
	CarId        *int
	WaitingSince time.Time
}

func (c EnqueueJourneyRequest) IsValid() bool {
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

	journeyID := req.Form.Get("ID")
	if journeyID == "" {
		return fmt.Errorf("ID is required")
	}

	r.ID, err = strconv.Atoi(journeyID)
	if err != nil {
		return fmt.Errorf("ID must be an integer")
	}

	if r.ID <= 0 {
		return fmt.Errorf("ID must be a positive integer")
	}

	return nil
}

type DropoffRequest struct {
	ID int
}

func (r *DropoffRequest) Validate(req *http.Request) error {
	err := req.ParseForm()

	if err != nil {
		return fmt.Errorf("bad request")
	}

	journeyID := req.Form.Get("ID")
	if journeyID == "" {
		return fmt.Errorf("ID is required")
	}

	r.ID, err = strconv.Atoi(journeyID)
	if err != nil {
		return fmt.Errorf("ID must be an integer")
	}

	if r.ID <= 0 {
		return fmt.Errorf("ID must be a positive integer")
	}

	return nil
}
