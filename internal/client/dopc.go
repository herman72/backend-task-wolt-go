package client

import (
	"backend-wolt-go/internal/models"
	"backend-wolt-go/internal/utils"
	"context"
)

type VenueProvider interface {
	GetVenueInformation(ctx context.Context, venueSlug string) (*models.VenueStaticResponse, *models.VenueDynamicResponse, error)
}

type DOPC struct {
	venueProvider VenueProvider
}

func NewDOPC(venueProvider VenueProvider) *DOPC {
	return &DOPC{venueProvider: venueProvider}
}

func (d *DOPC) CalculateDeliveryFee(ctx context.Context, orderInfo *models.OrderInfo) (models.PriceResponse, error) {
	staticResponse, dynamicResponse, err := d.venueProvider.GetVenueInformation(ctx, orderInfo.Slug)

	if err != nil {
		return models.PriceResponse{}, err
	}

	venueLon := staticResponse.VenueRaw.Location.Coordinates[0]
	venueLat := staticResponse.VenueRaw.Location.Coordinates[1]

	distance := utils.CalculateDistance(venueLat, venueLon, orderInfo.Lat, orderInfo.Lon)

	distanceRanges := make([]models.DistanceRange, len(dynamicResponse.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges))

	for i, dr := range dynamicResponse.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges {
		distanceRanges[i] = models.DistanceRange{
			Min: dr.Min,
			Max: dr.Max,
			A:   dr.A,
			B:   dr.B,
		}
	}

	deliveryFee, err := utils.CalculateDeliveryFee(distance, dynamicResponse.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice, distanceRanges)

	if err != nil {
		return models.PriceResponse{}, err
	}

	smallOrderSurcharge := utils.CalculateSmallOrderSurcharge(orderInfo.CartValue, dynamicResponse.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge)

	totalPrice := utils.CalculateTotalPrice(orderInfo.CartValue, smallOrderSurcharge, deliveryFee)

	return models.PriceResponse{
		TotalPrice: totalPrice,
		SmallOrderSurcharge: smallOrderSurcharge,
		CartValue: orderInfo.CartValue,
		Delivery: struct {
			Fee      int `json:"fee"`
			Distance int `json:"distance"`
		}{
			Fee:      deliveryFee,
			Distance: int(distance),
		},
	}, nil


}
