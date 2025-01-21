package service

import (
	"fmt"
	"math"
)

type CalculatorService interface {
	CalculateDistance(userLat, userLon, venueLat, venueLon float64) int
	CalculateDeliveryFee(distance, basePrice int, distanceRanges []DistanceRange) (int, error)
	CalculateSmallOrderSurcharge(cartValue, orderMinimumNoSurcharge int) int
	CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee int) int
}

type calculatorService struct {}

func NewCalculatorService() CalculatorService {
	return &calculatorService{}
}

// CalculateDistance computes the straight-line distance between two geographic points.
func (c *calculatorService) CalculateDistance(userLat, userLon, venueLat, venueLon float64) int {
	const earthRadius = 6371000 // Earth radius in meters
	latDiff := degreesToRadians(venueLat - userLat)
	lonDiff := degreesToRadians(venueLon - userLon)

	a := math.Sin(latDiff/2)*math.Sin(latDiff/2) +
		math.Cos(degreesToRadians(userLat))*math.Cos(degreesToRadians(venueLat))*
		math.Sin(lonDiff/2)*math.Sin(lonDiff/2)

	cVal := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * cVal
	return int(distance)
}

// CalculateDeliveryFee computes the delivery fee based on distance ranges and a base price.
func (c *calculatorService) CalculateDeliveryFee(distance, basePrice int, distanceRanges []DistanceRange) (int, error) {

	var fee int
	for _, rangeData := range distanceRanges {
		if distance >= rangeData.Min && distance < rangeData.Max {
			fee = basePrice + rangeData.A + int(rangeData.B*float64(distance)/10)
		}
	}
	if fee == 0 {
		return 0, fmt.Errorf("delivery is not possible, distance too long")
	}
	return fee, nil

}

// CalculateSmallOrderSurcharge computes the surcharge for small orders.
func (c *calculatorService) CalculateSmallOrderSurcharge(cartValue, orderMinimumNoSurcharge int) int {
	if cartValue >= orderMinimumNoSurcharge {
		return 0
	}
	return orderMinimumNoSurcharge - cartValue
}

// CalculateTotalPrice sums up the cart value, small order surcharge, and delivery fee.
func (c *calculatorService) CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee int) int {
	return cartValue + smallOrderSurcharge + deliveryFee
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

type DistanceRange struct {
	Min int     `json:"min"`
	Max int     `json:"max"`
	A   int     `json:"a"`
	B   float64 `json:"b"`
}




// func CalculateDeliveryFee(distance int, basePrice int, distanceRange []struct {
// 	Min int     `json:"min"`
// 	Max int     `json:"max"`
// 	A   int     `json:"a"`
// 	B   float64 `json:"b"`
// }) (int, error) {
// 	var fee int
// 	for _, rangeData := range distanceRange {
// 		if distance >= rangeData.Min && distance < rangeData.Max {
// 			fee = basePrice + rangeData.A + int(rangeData.B*float64(distance)/10)
// 		}
// 	}
// 	if fee == 0 {
// 		return 0, fmt.Errorf("delivery is not possible, distance too long")
// 	}
// 	return fee, nil
// }

// func CalculateSmallOrderSurcharge(cartValue int, orderMinimumNoSurcharge int) int {
	
// 	return int(math.Abs(float64(cartValue - orderMinimumNoSurcharge)))
// }

// func CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee int) int {
// 	return cartValue + smallOrderSurcharge + deliveryFee
// }