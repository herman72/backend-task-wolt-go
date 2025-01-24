package client

import (
	"backend-wolt-go/internal/models"
	"backend-wolt-go/internal/utils"
	"context"
)

// VenueProvider defines an interface for retrieving venue information.
type VenueProvider interface {
	// GetVenueInformation retrieves static and dynamic information for a given venue by its slug.
	GetVenueInformation(ctx context.Context, venueSlug string) (*models.VenueStaticResponse, *models.VenueDynamicResponse, error)
}

// DOPC (Delivery Order Price Calculator) is responsible for calculating delivery fees.
type DOPC struct {
	venueProvider VenueProvider
}

// NewDOPC creates a new instance of the DOPC struct with the provided VenueProvider.
func NewDOPC(venueProvider VenueProvider) *DOPC {
	return &DOPC{venueProvider: venueProvider}
}

// CalculateDeliveryFee calculates the delivery fee based on order information.
// It retrieves venue data, calculates the distance between the venue and the user,
// determines the delivery fee, small order surcharge, and total price.
func (d *DOPC) CalculateDeliveryFee(ctx context.Context, orderInfo *models.OrderInfo) (models.PriceResponse, error) {
	// Retrieve venue information (static and dynamic) for the given venue slug.
	staticResponse, dynamicResponse, err := d.venueProvider.GetVenueInformation(ctx, orderInfo.Slug)
	if err != nil {
		return models.PriceResponse{}, err
	}

	// Extract venue coordinates from the static response.
	venueLon := staticResponse.VenueRaw.Location.Coordinates[0]
	venueLat := staticResponse.VenueRaw.Location.Coordinates[1]

	// Calculate the distance between the venue and the user's location.
	distance := utils.CalculateDistance(venueLat, venueLon, orderInfo.Lat, orderInfo.Lon)

	// Map the distance ranges from the dynamic response to a usable format.
	distanceRanges := make([]models.DistanceRange, len(dynamicResponse.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges))
	for i, dr := range dynamicResponse.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges {
		distanceRanges[i] = models.DistanceRange{
			Min: dr.Min,
			Max: dr.Max,
			A:   dr.A,
			B:   dr.B,
		}
	}

	// Calculate the delivery fee based on the distance, base price, and distance ranges.
	deliveryFee, err := utils.CalculateDeliveryFee(distance, dynamicResponse.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice, distanceRanges)
	if err != nil {
		return models.PriceResponse{}, err
	}

	// Calculate the small order surcharge if the cart value is below the minimum threshold.
	smallOrderSurcharge := utils.CalculateSmallOrderSurcharge(orderInfo.CartValue, dynamicResponse.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge)

	// Calculate the total price, including cart value, small order surcharge, and delivery fee.
	totalPrice := utils.CalculateTotalPrice(orderInfo.CartValue, smallOrderSurcharge, deliveryFee)

	// Return the calculated price response.
	return models.PriceResponse{
		TotalPrice:          totalPrice,
		SmallOrderSurcharge: smallOrderSurcharge,
		CartValue:           orderInfo.CartValue,
		Delivery: struct {
			Fee      int `json:"fee"`
			Distance int `json:"distance"`
		}{
			Fee:      deliveryFee,
			Distance: int(distance),
		},
	}, nil
}
