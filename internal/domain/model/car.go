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

type AddCarRequest struct {
	ID    int `json:"id"`
	Seats int `json:"seats"`
}

func (c AddCarRequest) IsValid() bool {
	return c.Seats >= 4 && c.Seats <= 6
}

type CarResponse struct {
	ID    int `json:"id"`
	Seats int `json:"seats"`
}
