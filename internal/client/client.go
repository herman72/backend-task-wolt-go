package client

import (
	"backend-wolt-go/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type APIClient interface {
	GetVenueData(venueSlug string) (*models.VenueStaticResponse, *models.VenueDynamicResponse, error)
}
type apiClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewAPIClient() APIClient {
	return &apiClient{
		httpClient: &http.Client{},
		baseURL:    "https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues",
	}
}

func (c *apiClient) GetVenueData(venueSlug string) (*models.VenueStaticResponse, *models.VenueDynamicResponse, error) {
	staticURL := fmt.Sprintf("%s/%s/static", c.baseURL, venueSlug)
	dynamicURL := fmt.Sprintf("%s/%s/dynamic", c.baseURL, venueSlug)

	staticData := &models.VenueStaticResponse{}
	dynamicData := &models.VenueDynamicResponse{}

	staticData, err := c.fetchStaticResponse(staticURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch static data: %w", err)
	}

	// if err := c.fetchJSON(staticURL, staticData); err != nil {
	// 	return nil, nil, fmt.Errorf("failed to fetch static data: %w", err)
	// }

	if err := c.fetchJSON(dynamicURL, dynamicData); err != nil {
		return nil, nil, fmt.Errorf("failed to fetch dynamic data: %w", err)
	}
	return staticData, dynamicData, nil
}

func (c *apiClient) fetchJSON(url string, target interface{}) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return nil
}

func (c *apiClient) fetchStaticResponse(url string) (*models.VenueStaticResponse, error) {
	respByte, err := c.callAPI(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call static route: %w", err)
	}

	staticData := &models.VenueStaticResponse{}

	err = json.Unmarshal(respByte, staticData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if staticData.VenueRaw == nil {
		var errorResponse ServerError
		_ = json.Unmarshal(respByte, &errorResponse)
		return nil, fmt.Errorf("failed to get response: %w", errorResponse)
	}

	return staticData, nil
}

func (c *apiClient) callAPI(url string) ([]byte, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func FetchVenueStatic(venueSlug string) (models.VenueStaticResponse, error) {
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

func FetchVenueDynamic(venueSlug string) (models.VenueDynamicResponse, error) {
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

type ServerError struct {
	Msg string `json:"message"`
}

func (e ServerError) Error() string {
	return e.Msg
}
