package model

import (
	"time"
)

type Journey struct {
	ID           int
	People       int
	CarId        *int
	WaitingSince time.Time
}

func NewJourney(ID int, people int) *Journey {
	return &Journey{
		ID:           ID,
		People:       people,
		CarId:        nil,
		WaitingSince: time.Now(),
	}
}
