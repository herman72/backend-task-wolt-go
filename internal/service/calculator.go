package service

import (
	"fmt"
	"math"
)





func CalculateDeliveryFee(distance int, basePrice int, distanceRange []struct {
	Min int     `json:"min"`
	Max int     `json:"max"`
	A   int     `json:"a"`
	B   float64 `json:"b"`
}) (int, error) {
	var fee int
	for _, rangeData := range distanceRange {
		if distance >= rangeData.Min && distance < rangeData.Max {
			fee = basePrice + rangeData.A + int(rangeData.B*float64(distance)/10)
		}
	}
	if fee == 0 {
		return 0, fmt.Errorf("delivery is not possible, distance too long")
	}
	return fee, nil
}

func CalculateSmallOrderSurcharge(cartValue int, orderMinimumNoSurcharge int) int {
	
	return int(math.Abs(float64(cartValue - orderMinimumNoSurcharge)))
}

func CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee int) int {
	return cartValue + smallOrderSurcharge + deliveryFee
}