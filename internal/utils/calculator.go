package utils

import (
	"backend-wolt-go/internal/models"
	"fmt"
	"math"
)

// CalculateDeliveryFee calculates the delivery fee based on the distance,
// base price, and a range of distance-based pricing rules.
// Returns the delivery fee or an error if the distance exceeds the supported range.
func CalculateDeliveryFee(distance int, basePrice int, distanceRange []models.DistanceRange) (int, error) {
	var fee int
	var deliveryFeeFlag bool = false
	for _, rangeData := range distanceRange {
		if distance >= rangeData.Min && distance < rangeData.Max {
			fee = basePrice + rangeData.A + int(rangeData.B*float64(distance)/10)
			deliveryFeeFlag = true
			break
		}
	}
	if !deliveryFeeFlag {
		return 0, fmt.Errorf("distance exceeds allowable range, distance too long: %d meters", distance)
	}
	if deliveryFeeFlag && fee < 0 {
		return 0, fmt.Errorf("delivery fee cannot be negative value %d", fee)
	}
	return fee, nil
}

// CalculateSmallOrderSurcharge calculates the surcharge for small orders
// if the cart value is below the minimum value required to avoid the surcharge.
func CalculateSmallOrderSurcharge(cartValue int, orderMinimumNoSurcharge int) int {
	if cartValue < orderMinimumNoSurcharge {
		return orderMinimumNoSurcharge - cartValue
	}
	return 0
}

// CalculateTotalPrice calculates the total price of an order,
// including the cart value, small order surcharge, and delivery fee.
func CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee int) int {
	return cartValue + smallOrderSurcharge + deliveryFee
}

// CalculateDistance calculates the great-circle distance between two points
// on the Earth specified by their latitude and longitude.
// The distance is returned in meters.
func CalculateDistance(lat1, lon1, lat2, lon2 float64) int {
	const EarthRadius = 6371000 // Earth radius in meters

	// Convert latitudes and longitudes from degrees to radians.
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Compute differences in latitude and longitude.
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	// Apply the haversine formula to calculate the distance.
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int(math.Round(EarthRadius * c))
}
