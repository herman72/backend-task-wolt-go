package api

import (
	"backend-wolt-go/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

// DOPCService defines the interface for a service that calculates delivery fees.
type DOPCService interface {
	// CalculateDeliveryFee calculates the delivery fee based on order information.
	CalculateDeliveryFee(*models.OrderInfo) (models.PriceResponse, error)
}

// Handler is the HTTP handler for delivery order price calculation.
type Handler struct {
	service DOPCService
}

// NewHandler creates a new Handler instance with the provided DOPCService.
func NewHandler(service DOPCService) *Handler {
	return &Handler{service: service}
}

// GetDeliveryOrderPrice handles HTTP GET requests for calculating delivery order prices.
// It validates query parameters, constructs the order information, calls the service,
// and returns the delivery fee as a JSON response.
func (h *Handler) GetDeliveryOrderPrice(w http.ResponseWriter, r *http.Request) {
	// Extract and validate the venue_slug parameter.
	venueSlug := r.URL.Query().Get("venue_slug")
	if venueSlug == "" {
		http.Error(w, "Missing required parameter: venue_slug", http.StatusBadRequest)
		return
	}

	// Extract and validate the user_lat parameter.
	latStr := r.URL.Query().Get("user_lat")
	if latStr == "" {
		http.Error(w, "Missing required parameter: user_lat", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid user latitude", http.StatusBadRequest)
		return
	}

	if lat < -90 || lat > 90 {
		http.Error(w, "Latitude must be between -90 and 90", http.StatusBadRequest)
		return
	}

	// Extract and validate the user_lon parameter.
	lonStr := r.URL.Query().Get("user_lon")
	if lonStr == "" {
		http.Error(w, "Missing required parameter: user_lon", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Invalid user longitude", http.StatusBadRequest)
		return
	}

	if lon < -180 || lon > 180 {
		http.Error(w, "Longitude must be between -180 and 180", http.StatusBadRequest)
		return
	}

	// Extract and validate the cart_value parameter.
	cartValueStr := r.URL.Query().Get("cart_value")
	if cartValueStr == "" {
		http.Error(w, "Missing required parameter: cart_value", http.StatusBadRequest)
		return
	}

	cartValue, err := strconv.Atoi(cartValueStr)
	if err != nil {
		http.Error(w, "Invalid cart value", http.StatusBadRequest)
		return
	}

	if cartValue <= 0 {
		http.Error(w, "Cart value must be a positive integer", http.StatusBadRequest)
		return
	}

	// Create an OrderInfo struct with the validated parameters.
	orderInfo := &models.OrderInfo{
		Slug:      venueSlug,
		Lat:       lat,
		Lon:       lon,
		CartValue: cartValue,
	}

	// Call the service to calculate the delivery fee.
	response, err := h.service.CalculateDeliveryFee(orderInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the response content type to JSON and encode the response.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}