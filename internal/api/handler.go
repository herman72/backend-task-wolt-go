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