package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"backend-wolt-go/internal/models"
)

func FetchVenueStatic(venueSlug string)(models.VenueStaticResponse, error){
	var staticResponse models.VenueStaticResponse
	url := fmt.Sprintf("https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues/%s/static", venueSlug)
	resp, err := http.Get(url)
	if err != nil {
		return models.VenueStaticResponse{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return models.VenueStaticResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&staticResponse); err != nil {
		return staticResponse, err
	}
	return staticResponse, nil
}

func FetchVenueDynamic(venueSlug string)(models.VenueDynamicResponse, error){
	var dynamicResponse models.VenueDynamicResponse
	url := fmt.Sprintf("https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues/%s/dynamic", venueSlug)
	resp, err := http.Get(url)
	if err != nil {
		return models.VenueDynamicResponse{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return models.VenueDynamicResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&dynamicResponse); err != nil {
		return dynamicResponse, err
	}
	
	return dynamicResponse, nil
}