package model

type Car struct {
	ID      int
	Seats   int
	InQueue int
}

func NewCar(ID int, seats int) *Car {
	return &Car{
		ID:      ID,
		Seats:   seats,
		InQueue: seats, // Initially, the car will be in the queue corresponding to its total seat count
	}
}
