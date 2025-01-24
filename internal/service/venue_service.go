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

// VenueProvider is responsible for fetching static and dynamic venue information from a remote server.
type VenueProvider struct {
	client  *http.Client // HTTP client for making API calls.
	baseURL string       // Base URL of the venue information API.
}

// NewVenueProvider creates a new instance of VenueProvider with the given base URL.
func NewVenueProvider(baseURL string) *VenueProvider {
	return &VenueProvider{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

// GetVenueInformation retrieves both static and dynamic information for a specific venue.
// It makes two API calls: one for static data and another for dynamic data.
func (v *VenueProvider) GetVenueInformation(ctx context.Context, venueSlug string) (*models.VenueStaticResponse, *models.VenueDynamicResponse, error) {
	// Construct API URLs for static and dynamic data.
	staticURL := fmt.Sprintf("%s/%s/static", v.baseURL, venueSlug)
	dynamicURL := fmt.Sprintf("%s/%s/dynamic", v.baseURL, venueSlug)

	// Fetch static data.
	staticData, err := v.FetchVenueStaticData(ctx, staticURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get static data: %w", err)
	}

	// Fetch dynamic data.
	dynamicData, err := v.FetchVenueDynamicData(ctx, dynamicURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get dynamic data: %w", err)
	}

	return staticData, dynamicData, nil
}

// FetchVenueStaticData fetches the static information of a venue from the given URL.
func (v *VenueProvider) FetchVenueStaticData(ctx context.Context, url string) (*models.VenueStaticResponse, error) {
	// Call the API and get the response bytes.
	respByte, err := v.callAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to call static route: %w", err)
	}

	// Unmarshal the response into VenueStaticResponse.
	staticDataResponse := &models.VenueStaticResponse{}
	err = json.Unmarshal(respByte, staticDataResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Check if the response contains valid data.
	if staticDataResponse.VenueRaw == nil {
		var errorResponse models.ServerError
		_ = json.Unmarshal(respByte, &errorResponse)
		return nil, fmt.Errorf("failed to get response: %v", errorResponse)
	}

	return staticDataResponse, nil
}

// FetchVenueDynamicData fetches the dynamic information of a venue from the given URL.
func (v *VenueProvider) FetchVenueDynamicData(ctx context.Context, url string) (*models.VenueDynamicResponse, error) {
	// Call the API and get the response bytes.
	respByte, err := v.callAPI(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to call dynamic route: %w", err)
	}

	// Unmarshal the response into VenueDynamicResponse.
	dynamicDataResponse := &models.VenueDynamicResponse{}
	err = json.Unmarshal(respByte, dynamicDataResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Check if the response contains valid data.
	if dynamicDataResponse.VenueRaw == nil {
		var errorResponse models.ServerError
		_ = json.Unmarshal(respByte, &errorResponse)
		return nil, fmt.Errorf("failed to get response: %v", errorResponse)
	}

	return dynamicDataResponse, nil
}

// callAPI makes an HTTP GET request to the given URL and returns the response body as a byte slice.
func (v *VenueProvider) callAPI(ctx context.Context, url string) ([]byte, error) {
	// Create a new HTTP request with the provided context and URL.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the HTTP request.
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success.
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	// Read and return the response body.
	return io.ReadAll(resp.Body)
}
