package api

import (
	"backend-wolt-go/internal/models"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type DOPCService interface {
	CalculateDeliveryFee(context.Context, *models.OrderInfo) (models.PriceResponse, error)
}

type Handler struct {
	service DOPCService
}

func NewHandler(service DOPCService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetDeliveryOrderPrice(w http.ResponseWriter, r *http.Request) {
	venueSlug := r.URL.Query().Get("venue_slug")
	if venueSlug == "" {
		http.Error(w, "Missing required parameter: venue_slug", http.StatusBadRequest)
		return
	}

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

	orderInfo := &models.OrderInfo{
		Slug:      venueSlug,
		Lat:       lat,
		Lon:       lon,
		CartValue: cartValue,
	}

	response, err := h.service.CalculateDeliveryFee(r.Context(), orderInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}


}