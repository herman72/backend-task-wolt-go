package api

import (
	"backend-wolt-go/internal/client"
	"fmt"
	"net/http"
	// "strconv"
)


func Handler(w http.ResponseWriter, r *http.Request){
	venueSlug := r.URL.Query().Get("venue_slug")
	cartValueStr := r.URL.Query().Get("cart_value")
	userLatStr := r.URL.Query().Get("user_lat")
	userLonStr := r.URL.Query().Get("user_lon")

	if venueSlug == "" || cartValueStr == "" || userLatStr == "" || userLonStr == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// cartValue, err := strconv.Atoi(cartValueStr)
	// if err != nil {
	// 	http.Error(w, "Invalid cart value", http.StatusBadRequest)
	// }

	// userLat, err := strconv.ParseFloat(userLatStr, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid user latitude", http.StatusBadRequest)
	// }

	// userLon, err := strconv.ParseFloat(userLonStr, 64)
	// if err != nil {
	// 	http.Error(w, "Invalid user longitude", http.StatusBadRequest)
	// }

	staticData, err := client.FetchVenueStatic(venueSlug)
	if err != nil {
		http.Error(w, "Failed to fetch venue static data", http.StatusInternalServerError)
	}

	// dynamicData, err := client.FetchVenueDynamic(venueSlug)
	// if err != nil {
	// 	http.Error(w, "Failed to fetch venue dynamic data", http.StatusInternalServerError)
	// }

	fmt.Println(staticData)

	return

}

