package service

import (
	"backend-wolt-go/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)



type VenueProvider struct {
	client *http.Client
	baseURL string
}

func NewVenueProvider(baseURL string) *VenueProvider {
	return &VenueProvider{
		client: &http.Client{},
		baseURL: baseURL,
	}
}

func (v *VenueProvider) GetVenueInformation(ctx context.Context ,venueSlug string) (*models.VenueStaticResponse, *models.VenueDynamicResponse, error) {
	staticURL := fmt.Sprintf("%s/%s/static", v.baseURL, venueSlug)
	dynamicURL := fmt.Sprintf("%s/%s/dynamic", v.baseURL, venueSlug)

	staticData, err := v.FetchVenueStaticData(ctx, staticURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get static data: %w", err)
	}
	
	dynamicData, err := v.FetchVenueDynamicData(ctx, dynamicURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get dynamic data: %w", err)
	}

	return staticData, dynamicData, nil
}

func (v *VenueProvider) FetchVenueStaticData(ctx context.Context, url string) (*models.VenueStaticResponse, error){
	respByte, err := v.callAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to call static route: %w", err)
	}

	staticDataResponse := &models.VenueStaticResponse{}
	err = json.Unmarshal(respByte, staticDataResponse)

	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if staticDataResponse.VenueRaw == nil {
		var errorResponse models.ServerError
		_ = json.Unmarshal(respByte, &errorResponse)
		return nil, fmt.Errorf("failed to get response: %v", errorResponse)
	}

	return staticDataResponse, nil
	
}

func (v *VenueProvider) FetchVenueDynamicData(ctx context.Context, url string) (*models.VenueDynamicResponse, error){
	respByte, err := v.callAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to call dynamic route: %w", err)
	}

	dynamicDataResponse := &models.VenueDynamicResponse{}
	err = json.Unmarshal(respByte, dynamicDataResponse)

	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if dynamicDataResponse.VenueRaw == nil {
		var errorResponse models.ServerError
		_ = json.Unmarshal(respByte, &errorResponse)
		return nil, fmt.Errorf("failed to get response: %v", errorResponse)
	}

	return dynamicDataResponse, nil
}

func (v *VenueProvider) callAPI(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	return io.ReadAll(resp.Body)
}