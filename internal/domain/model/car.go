package model

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

type CarResponse struct {
	ID             int `json:"id"`
	Seats          int `json:"seats"`
	AvailableSeats int `json:"availableSeats"`
}
