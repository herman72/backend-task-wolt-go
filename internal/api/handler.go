package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend-wolt-go/internal/client"
	"backend-wolt-go/internal/service"
	"backend-wolt-go/internal/models"
)

type Handler struct {
	apiClient   client.APIClient
	calculator service.CalculatorService
}

func NewHandler(apiClient client.APIClient, calculator service.CalculatorService) *Handler {
	return &Handler{
		apiClient:   apiClient,
		calculator: calculator,
	}
}

func (h *Handler) GetDeliveryOrderPrice(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	venueSlug := r.URL.Query().Get("venue_slug")
	cartValueStr := r.URL.Query().Get("cart_value")
	userLatStr := r.URL.Query().Get("user_lat")
	userLonStr := r.URL.Query().Get("user_lon")

	if venueSlug == "" || cartValueStr == "" || userLatStr == "" || userLonStr == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	cartValue, err := strconv.Atoi(cartValueStr)
	if err != nil {
		http.Error(w, "Invalid cart value", http.StatusBadRequest)
		return
	}

	userLat, err := strconv.ParseFloat(userLatStr, 64)
	if err != nil {
		http.Error(w, "Invalid user latitude", http.StatusBadRequest)
		return
	}

	userLon, err := strconv.ParseFloat(userLonStr, 64)
	if err != nil {
		http.Error(w, "Invalid user longitude", http.StatusBadRequest)
		return
	}

	// Fetch venue data
	staticData, dynamicData, err := h.apiClient.GetVenueData(venueSlug)
	if err != nil {
		http.Error(w, "Failed to fetch venue data", http.StatusInternalServerError)
		return
	}

	vanueLon := staticData.VenueRaw.Location.Coordinates[0]
	vanueLat := staticData.VenueRaw.Location.Coordinates[1]

	// Calculate delivery distance
	distance := h.calculator.CalculateDistance(userLat, userLon, vanueLon, vanueLat)

	// Calculate delivery fee
	distanceRanges := make([]service.DistanceRange, len(dynamicData.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges))
	for i, dr := range dynamicData.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges {
		distanceRanges[i] = service.DistanceRange{
			Min: dr.Min,
			Max: dr.Max,
			A:   dr.A,
			B:   dr.B,
		}
	}

	deliveryFee, err := h.calculator.CalculateDeliveryFee(distance, dynamicData.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice, distanceRanges)
	if err != nil {
		http.Error(w, "Delivery not possible", http.StatusBadRequest)
		return
	}

	// Calculate small order surcharge
	smallOrderSurcharge := h.calculator.CalculateSmallOrderSurcharge(cartValue, dynamicData.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge)

	// Calculate total price
	totalPrice := h.calculator.CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee)

	// Create response
	response := models.PriceResponse{
		TotalPrice:          totalPrice,
		SmallOrderSurcharge: smallOrderSurcharge,
		CartValue:           cartValue,
		Delivery: struct {
			Fee      int `json:"fee"`
			Distance int `json:"distance"`
		}{
			Fee:      deliveryFee,
			Distance: int(distance),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}





// package api

// import (
// 	"backend-wolt-go/internal/client"
// 	"backend-wolt-go/internal/service"
// 	"backend-wolt-go/internal/utils"
// 	"backend-wolt-go/internal/models"
// 	"encoding/json"
// 	"net/http"
// 	"strconv"
// )


// func Handler(w http.ResponseWriter, r *http.Request){
// 	venueSlug := r.URL.Query().Get("venue_slug")
// 	cartValueStr := r.URL.Query().Get("cart_value")
// 	userLatStr := r.URL.Query().Get("user_lat")
// 	userLonStr := r.URL.Query().Get("user_lon")

// 	if venueSlug == "" || cartValueStr == "" || userLatStr == "" || userLonStr == "" {
// 		http.Error(w, "Missing required parameters", http.StatusBadRequest)
// 		return
// 	}

// 	cartValue, err := strconv.Atoi(cartValueStr)
// 	if err != nil {
// 		http.Error(w, "Invalid cart value", http.StatusBadRequest)
// 	}

// 	userLat, err := strconv.ParseFloat(userLatStr, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid user latitude", http.StatusBadRequest)
// 	}

// 	userLon, err := strconv.ParseFloat(userLonStr, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid user longitude", http.StatusBadRequest)
// 	}

// 	staticData, err := client.FetchVenueStatic(venueSlug)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch venue static data", http.StatusInternalServerError)
// 	}

// 	dynamicData, err := client.FetchVenueDynamic(venueSlug)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch venue dynamic data", http.StatusInternalServerError)
// 	}
// 	 vanueLon := staticData.VenueRaw.Location.Coordinates[0]
// 	 vanueLat := staticData.VenueRaw.Location.Coordinates[1]

// 	distance := utils.CalculateDistance(userLat, userLon, vanueLat, vanueLon)
	
// 	deliveryFee, err := service.CalculateDeliveryFee(int(distance), dynamicData.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice,
// 		dynamicData.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges)
// 	if err != nil {
// 		http.Error(w, "delivery not possible for the given distance", http.StatusBadRequest)
// 		return
// 	}
// 	smallOrderSurcharge := service.CalculateSmallOrderSurcharge(cartValue, dynamicData.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge)

// 	totalPrice := service.CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee)

// 	response := models.PriceResponse{
// 		TotalPrice:          totalPrice,
// 		SmallOrderSurcharge: smallOrderSurcharge,
// 		CartValue:           cartValue,
// 		Delivery: struct {
// 			Fee      int `json:"fee"`
// 			Distance int `json:"distance"`
// 		}{
// 			Fee:      deliveryFee,
// 			Distance: int(distance),
// 		},
// 	}

// 	// Send response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)

// }

