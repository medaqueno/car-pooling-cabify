package dto

import (
	"car-pooling-service/internal/errors"
	"net/http"
	"strconv"
)

type AddCarRequest struct {
	ID    int `json:"id"`
	Seats int `json:"seats"`
}

func (c AddCarRequest) IsValid() bool {
	return c.Seats >= 4 && c.Seats <= 6
}

type EnqueueJourneyRequest struct {
	ID     int `json:"id"`
	People int `json:"people"`
}

func (c EnqueueJourneyRequest) IsValid() bool {
	return c.People >= 1 && c.People <= 6
}

type LocateJourneyRequest struct {
	ID int
}

func (r *LocateJourneyRequest) Validate(req *http.Request) error {
	err := req.ParseForm()
	if err != nil {
		return errors.ErrGenericBadReq
	}

	journeyID := req.Form.Get("ID")
	if journeyID == "" {
		return errors.ErrIDRequired
	}

	r.ID, err = strconv.Atoi(journeyID)
	if err != nil {
		return errors.ErrInvalidID
	}

	if r.ID <= 0 {
		return errors.ErrIDPositive
	}

	return nil
}

type DropoffRequest struct {
	ID int
}

func (r *DropoffRequest) Validate(req *http.Request) error {
	err := req.ParseForm()
	if err != nil {
		return errors.ErrGenericBadReq
	}

	journeyID := req.Form.Get("ID")
	if journeyID == "" {
		return errors.ErrIDRequired
	}

	r.ID, err = strconv.Atoi(journeyID)
	if err != nil {
		return errors.ErrInvalidID
	}

	if r.ID <= 0 {
		return errors.ErrIDPositive
	}

	return nil
}
