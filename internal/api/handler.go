package api

import (
	"backend-wolt-go/internal/client"
	"backend-wolt-go/internal/service"
	"backend-wolt-go/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type PriceResponse struct {
	TotalPrice          int `json:"total_price"`
	SmallOrderSurcharge int `json:"small_order_surcharge"`
	CartValue           int `json:"cart_value"`
	Delivery            struct {
		Fee      int `json:"fee"`
		Distance int `json:"distance"`
	} `json:"delivery"`
}


func Handler(w http.ResponseWriter, r *http.Request){
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
	}

	userLat, err := strconv.ParseFloat(userLatStr, 64)
	if err != nil {
		http.Error(w, "Invalid user latitude", http.StatusBadRequest)
	}

	userLon, err := strconv.ParseFloat(userLonStr, 64)
	if err != nil {
		http.Error(w, "Invalid user longitude", http.StatusBadRequest)
	}

	staticData, err := client.FetchVenueStatic(venueSlug)
	if err != nil {
		http.Error(w, "Failed to fetch venue static data", http.StatusInternalServerError)
	}

	dynamicData, err := client.FetchVenueDynamic(venueSlug)
	if err != nil {
		http.Error(w, "Failed to fetch venue dynamic data", http.StatusInternalServerError)
	}
	 vanueLon := staticData.VenueRaw.Location.Coordinates[0]
	 vanueLat := staticData.VenueRaw.Location.Coordinates[1]

	distance := utils.CalculateDistance(userLat, userLon, vanueLat, vanueLon)
	
	deliveryFee, err := service.CalculateDeliveryFee(int(distance), dynamicData.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice,
		dynamicData.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges)
	if err != nil {
		http.Error(w, "delivery not possible for the given distance", http.StatusBadRequest)
		return
	}
	smallOrderSurcharge := service.CalculateSmallOrderSurcharge(cartValue, dynamicData.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge)

	totalPrice := service.CalculateTotalPrice(cartValue, smallOrderSurcharge, deliveryFee)

	response := PriceResponse{
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

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

