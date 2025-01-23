package client

import (
	
	"backend-wolt-go/internal/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ------------------------------------------------------------
// 1. Define a mock VenueProvider for testing
// ------------------------------------------------------------
type mockVenueProvider struct {
	mock.Mock
}

func (m *mockVenueProvider) GetVenueInformation(
	ctx context.Context,
	venueSlug string,
) (*models.VenueStaticResponse, *models.VenueDynamicResponse, error) {

	args := m.Called(ctx, venueSlug)

	// If you want to handle potential nil returns, you could do type-checking:
	staticResp, _ := args.Get(0).(*models.VenueStaticResponse)
	dynamicResp, _ := args.Get(1).(*models.VenueDynamicResponse)

	return staticResp, dynamicResp, args.Error(2)
}

// ------------------------------------------------------------
// 2. Write the test for CalculateDeliveryFee
// ------------------------------------------------------------
func TestDOPC_CalculateDeliveryFee(t *testing.T) {
	// Create our mock VenueProvider
	mockProvider := new(mockVenueProvider)

	// Instantiate DOPC with the mock
	dopc := NewDOPC(mockProvider)

	// Define some test input
	orderInfo := &models.OrderInfo{
		Slug:      "my-test-venue",
		Lat:       60.1708,  // user/customer latitude
		Lon:       24.9375,  // user/customer longitude
		CartValue: 1000,     // e.g. 10 euros in cents
	}

	// Define what static venue data the mock should return
	staticResp := &models.VenueStaticResponse{
		VenueRaw: &struct {
			Location struct {
				Coordinates []float64 `json:"coordinates"`
			} `json:"location"`
		}{
			Location: struct {
				Coordinates []float64 `json:"coordinates"`
			}{
				// Typically: Coordinates[0] = longitude, Coordinates[1] = latitude
				Coordinates: []float64{24.93545, 60.16952},
			},
		},
	}

	// Define what dynamic venue data the mock should return
	dynamicResp := &models.VenueDynamicResponse{
		VenueRaw: &struct {
			DeliverySpecs struct {
				OrderMinimumNoSurcharge int `json:"order_minimum_no_surcharge"`
				DeliveryPricing         struct {
					BasePrice      int `json:"base_price"`
					DistanceRanges []struct {
						Min int     `json:"min"`
						Max int     `json:"max"`
						A   int     `json:"a"`
						B   float64 `json:"b"`
					} `json:"distance_ranges"`
				} `json:"delivery_pricing"`
			} `json:"delivery_specs"`
		}{
			DeliverySpecs: struct {
				OrderMinimumNoSurcharge int `json:"order_minimum_no_surcharge"`
				DeliveryPricing         struct {
					BasePrice      int `json:"base_price"`
					DistanceRanges []struct {
						Min int     `json:"min"`
						Max int     `json:"max"`
						A   int     `json:"a"`
						B   float64 `json:"b"`
					} `json:"distance_ranges"`
				} `json:"delivery_pricing"`
			}{
				OrderMinimumNoSurcharge: 1000, // No surcharge if cart >= 1000
				DeliveryPricing: struct {
					BasePrice      int `json:"base_price"`
					DistanceRanges []struct {
						Min int     `json:"min"`
						Max int     `json:"max"`
						A   int     `json:"a"`
						B   float64 `json:"b"`
					} `json:"distance_ranges"`
				}{
					BasePrice: 300, // e.g. 3 euros in cents
					DistanceRanges: []struct {
						Min int     `json:"min"`
						Max int     `json:"max"`
						A   int     `json:"a"`
						B   float64 `json:"b"`
					}{
						// first range
						{
							Min: 0,
							Max: 2000,
							A:   100,
							B:   0.05,
						},
						// second range
						{
							Min: 2001,
							Max: 5000,
							A:   200,
							B:   0.06,
						},
					},
				},
			},
		},
	}

	// 3. Setup the mock to return the above static/dynamic responses
	mockProvider.
		On("GetVenueInformation", mock.Anything, "my-test-venue").
		Return(staticResp, dynamicResp, nil)

	// 4. Call the method under test
	result, err := dopc.CalculateDeliveryFee(context.Background(), orderInfo)

	// ------------------------------------------------------------
	// 5. Assertions
	// ------------------------------------------------------------
	assert.NoError(t, err, "CalculateDeliveryFee should not return an error")

	// Ensure the mock was called with the expected parameters
	mockProvider.AssertCalled(t, "GetVenueInformation", mock.Anything, "my-test-venue")
	mockProvider.AssertExpectations(t)

	// Basic checks on the returned struct
	assert.Equal(t, orderInfo.CartValue, result.CartValue, "CartValue should match")
	// If cartValue == 1000 and the orderMinimumNoSurcharge == 1000,
	// we expect smallOrderSurcharge to be 0.
	assert.Equal(t, 0, result.SmallOrderSurcharge, "Surcharge should be zero for cartValue >= min")

	// The total price will include the cart value + any delivery fee
	assert.True(t, result.TotalPrice >= orderInfo.CartValue, "Total price must be at least the cart value")

	// Check that the Fee is at least the base price (300). The actual amount
	// depends on the `utils.CalculateDeliveryFee` logic.  
	assert.True(t, result.Delivery.Fee >= 300, "Delivery fee should be >= base price")

	// Because we used mock data that places the venue ~some distance
	// from the user, we at least verify distance is > 0
	assert.True(t, result.Delivery.Distance > 0, "Distance should be greater than 0")
}
