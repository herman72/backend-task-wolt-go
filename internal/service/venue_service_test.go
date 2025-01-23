package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetVenueInformation_Success(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.String(), "static") {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"venue_raw": {"location": {"coordinates": [24.9354, 60.1699]}}}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"venue_raw": {"delivery_specs": {"order_minimum_no_surcharge": 15, "delivery_pricing": {"base_price": 5}}}}`))
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	venueProvider := NewVenueProvider(server.URL)
	staticResp, dynamicResp, err := venueProvider.GetVenueInformation(context.Background(), "test-slug")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if staticResp.VenueRaw.Location.Coordinates[0] != 24.9354 {
		t.Errorf("unexpected coordinates: %v", staticResp.VenueRaw.Location.Coordinates)
	}
	if dynamicResp.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge != 15 {
		t.Errorf("unexpected order minimum: %d", dynamicResp.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge)
	}
}

func TestGetVenueInformation_Error(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	venueProvider := NewVenueProvider(server.URL)
	_, _, err := venueProvider.GetVenueInformation(context.Background(), "test-slug")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
