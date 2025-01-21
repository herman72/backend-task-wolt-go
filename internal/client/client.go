package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type VenueStaticResponse struct {
	VenueRaw struct {
		Location struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"location"`
	} `json:"venue_raw"`
}

type VenueDynamicResponse struct {
	VenueRaw struct {
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
	} `json:"venue_raw"`
}

func FetchVenueStatic(venueSlug string)(VenueStaticResponse, error){
	var staticResponse VenueStaticResponse
	url := fmt.Sprintf("https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues/%s/static", venueSlug)
	resp, err := http.Get(url)
	if err != nil {
		return VenueStaticResponse{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return VenueStaticResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&staticResponse); err != nil {
		return staticResponse, err
	}
	return staticResponse, nil
}

func FetchVenueDynamic(venueSlug string)(VenueDynamicResponse, error){
	var dynamicResponse VenueDynamicResponse
	url := fmt.Sprintf("https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues/%s/dynamic", venueSlug)
	resp, err := http.Get(url)
	if err != nil {
		return VenueDynamicResponse{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return VenueDynamicResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&dynamicResponse); err != nil {
		return dynamicResponse, err
	}
	
	return dynamicResponse, nil
}