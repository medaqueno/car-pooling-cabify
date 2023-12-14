package command

import (
	"car-pooling-service/internal/domain/repository"
	"car-pooling-service/internal/port/http/dto"
)

type DropoffHandler struct {
	repo repository.CarAssignerRepository
}

func NewDropoffHandler(repo repository.CarAssignerRepository) *DropoffHandler {
	return &DropoffHandler{
		repo: repo,
	}
}

func (h *DropoffHandler) Handle(dropoffRequest dto.DropoffRequest) error {
	err := h.repo.DequeueJourney(dropoffRequest.ID)

	return err
}
