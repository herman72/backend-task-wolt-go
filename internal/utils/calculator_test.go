package utils

import (
	"backend-wolt-go/internal/models"
	"testing"
)

func TestCalculateDeliveryFee_Success(t *testing.T) {
	distance := 15
	basePrice := 10
	distanceRanges := []models.DistanceRange{
		{Min: 0, Max: 10, A: 5, B: 1},
		{Min: 10, Max: 20, A: 10, B: 2},
	}
	fee, err := CalculateDeliveryFee(distance, basePrice, distanceRanges)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if fee != 23 {
		t.Errorf("expected fee to be 23, got %d", fee)
	}
}

func TestCalculateDeliveryFee_OutOfRange(t *testing.T) {
	distance := 50
	basePrice := 10
	distanceRanges := []models.DistanceRange{
		{Min: 0, Max: 10, A: 5, B: 1},
		{Min: 10, Max: 20, A: 10, B: 2},
	}
	_, err := CalculateDeliveryFee(distance, basePrice, distanceRanges)
	if err == nil {
		t.Error("expected error for out of range distance, got nil")
	}
}