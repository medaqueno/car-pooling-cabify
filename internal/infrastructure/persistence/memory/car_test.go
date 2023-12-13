package memory

import (
	"car-pooling-service/internal/domain/model"
	"testing"
)

func TestAddCar(t *testing.T) {
	repo := NewCarRepository()
	car := model.NewCar(1, 4)

	if err := repo.AddCar(car); err != nil {
		t.Errorf("AddCar failed: %v", err)
	}

	// Test adding a duplicate car
	if err := repo.AddCar(car); err == nil {
		t.Error("Expected an error when adding a duplicate car, but got none")
	}
}

func TestFindCarByID(t *testing.T) {
	repo := NewCarRepository()
	car := model.NewCar(1, 4)
	_ = repo.AddCar(car)

	if _, err := repo.FindCarByID(1); err != nil {
		t.Errorf("FindCarByID failed to find car: %v", err)
	}

	// Test finding a non-existing car
	if _, err := repo.FindCarByID(2); err == nil {
		t.Error("Expected an error when finding a non-existing car, but got none")
	}
}

func TestUpdateCar(t *testing.T) {
	repo := NewCarRepository()
	car := model.NewCar(1, 4)
	_ = repo.AddCar(car)

	car.Seats = 5
	if err := repo.UpdateCar(car); err != nil {
		t.Errorf("UpdateCar failed: %v", err)
	}

	// Test updating a non-existing car
	nonExistingCar := model.NewCar(2, 4)
	if err := repo.UpdateCar(nonExistingCar); err == nil {
		t.Error("Expected an error when updating a non-existing car, but got none")
	}
}
