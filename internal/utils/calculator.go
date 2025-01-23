package utils

import (
	"backend-wolt-go/internal/models"
	"fmt"
	"math"
)

func CalculateDeliveryFee(distance int, basePrice int, distanceRange []models.DistanceRange) (int, error) {
	var fee int
	for _, rangeData := range distanceRange {
		if distance >= rangeData.Min && distance < rangeData.Max {
			fee = basePrice + rangeData.A + int(rangeData.B*float64(distance)/10)
			break // Stop looping once the correct range is found
		}
	}
	if fee == 0 {
		return 0, fmt.Errorf("delivery is not possible, distance too long")
	}
	return fee, nil
}

func CalculateSmallOrderSurcharge(cartValue int, orderMinimumNoSurcharge int) int {
	if cartValue < orderMinimumNoSurcharge {
		return orderMinimumNoSurcharge - cartValue
	}
	return 0
}

func CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee int) int {
	return cartValue + smallOrderSurcharge + deliveryFee
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) int {
	const EarthRadius = 6371000 // Earth radius in meters
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int(math.Round(EarthRadius * c))
}