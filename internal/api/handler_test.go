package api

import (
	
	"backend-wolt-go/internal/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ------------------------------
// 1. Mock the DOPCService
// ------------------------------
type mockDOPCService struct {
	mock.Mock
}

func (m *mockDOPCService) CalculateDeliveryFee(ctx context.Context, orderInfo *models.OrderInfo) (models.PriceResponse, error) {
	args := m.Called(ctx, orderInfo)
	resp, _ := args.Get(0).(models.PriceResponse)
	return resp, args.Error(1)
}

// ------------------------------
// 2. Helper to build request
// ------------------------------
func buildRequest(params map[string]string) *http.Request {
	// Build query string
	q := make([]string, 0, len(params))
	for k, v := range params {
		q = append(q, k+"="+v)
	}
	url := "/delivery-order-price?" + strings.Join(q, "&")

	req := httptest.NewRequest(http.MethodGet, url, nil)
	return req
}

// ------------------------------
// 3. Test missing/invalid params
// ------------------------------
func TestGetDeliveryOrderPrice_BadRequest(t *testing.T) {
	service := new(mockDOPCService)
	handler := NewHandler(service)

	tests := []struct {
		name       string
		params     map[string]string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Missing venue_slug",
			params:     map[string]string{"user_lat": "60.1699", "user_lon": "24.9384", "cart_value": "1500"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Missing required parameter: venue_slug",
		},
		{
			name:       "Missing user_lat",
			params:     map[string]string{"venue_slug": "abc", "user_lon": "24.9384", "cart_value": "1500"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Missing required parameter: user_lat",
		},
		{
			name:       "Invalid user_lat",
			params:     map[string]string{"venue_slug": "abc", "user_lat": "invalid", "user_lon": "24.9384", "cart_value": "1500"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid user latitude",
		},
		{
			name:       "Out of range user_lat",
			params:     map[string]string{"venue_slug": "abc", "user_lat": "100", "user_lon": "24.9384", "cart_value": "1500"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Latitude must be between -90 and 90",
		},
		{
			name:       "Missing user_lon",
			params:     map[string]string{"venue_slug": "abc", "user_lat": "60.1699", "cart_value": "1500"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Missing required parameter: user_lon",
		},
		{
			name:       "Invalid user_lon",
			params:     map[string]string{"venue_slug": "abc", "user_lat": "60.1699", "user_lon": "invalid", "cart_value": "1500"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid user longitude",
		},
		{
			name:       "Out of range user_lon",
			params:     map[string]string{"venue_slug": "abc", "user_lat": "60.1699", "user_lon": "200", "cart_value": "1500"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Longitude must be between -180 and 180",
		},
		{
			name:       "Missing cart_value",
			params:     map[string]string{"venue_slug": "abc", "user_lat": "60.1699", "user_lon": "24.9384"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Missing required parameter: cart_value",
		},
		{
			name:       "Invalid cart_value",
			params:     map[string]string{"venue_slug": "abc", "user_lat": "60.1699", "user_lon": "24.9384", "cart_value": "invalid"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid cart value",
		},
		{
			name:       "Zero or negative cart_value",
			params:     map[string]string{"venue_slug": "abc", "user_lat": "60.1699", "user_lon": "24.9384", "cart_value": "0"},
			wantStatus: http.StatusBadRequest,
			wantBody:   "Cart value must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := buildRequest(tt.params)

			handler.GetDeliveryOrderPrice(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.wantBody)
		})
	}
}

// ------------------------------
// 4. Test successful scenario
// ------------------------------
func TestGetDeliveryOrderPrice_Success(t *testing.T) {
	service := new(mockDOPCService)
	handler := NewHandler(service)

	// Expected result from our mock
	expectedPriceResponse := models.PriceResponse{
		TotalPrice:          2350,
		SmallOrderSurcharge: 0,
		CartValue:           2000,
		Delivery: struct {
			Fee      int `json:"fee"`
			Distance int `json:"distance"`
		}{
			Fee:      350,
			Distance: 987,
		},
	}

	// Set up the mock to expect certain input and return `expectedPriceResponse`
	service.On(
		"CalculateDeliveryFee",
		mock.Anything,
		&models.OrderInfo{
			Slug:      "venue123",
			Lat:       60.1699,
			Lon:       24.9384,
			CartValue: 2000,
		},
	).Return(expectedPriceResponse, nil)

	params := map[string]string{
		"venue_slug": "venue123",
		"user_lat":   "60.1699",
		"user_lon":   "24.9384",
		"cart_value": "2000",
	}

	rec := httptest.NewRecorder()
	req := buildRequest(params)

	handler.GetDeliveryOrderPrice(rec, req)

	// Validate the status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Validate the JSON body
	var actualResp models.PriceResponse
	err := json.Unmarshal(rec.Body.Bytes(), &actualResp)
	assert.NoError(t, err)

	assert.Equal(t, expectedPriceResponse, actualResp)

	// Verify that our mock was called with the right arguments
	service.AssertExpectations(t)
}

// ------------------------------
// 5. Test service error scenario
// ------------------------------
func TestGetDeliveryOrderPrice_ServiceError(t *testing.T) {
	service := new(mockDOPCService)
	handler := NewHandler(service)

	serviceErr := errors.New("some internal error")
	service.On(
		"CalculateDeliveryFee",
		mock.Anything,
		mock.AnythingOfType("*models.OrderInfo"),
	).Return(models.PriceResponse{}, serviceErr)

	params := map[string]string{
		"venue_slug": "venue-slug",
		"user_lat":   "60.1699",
		"user_lon":   "24.9384",
		"cart_value": "1500",
	}

	rec := httptest.NewRecorder()
	req := buildRequest(params)

	handler.GetDeliveryOrderPrice(rec, req)

	// Expect 500 status for service errors
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), serviceErr.Error())

	// Verify the mock was called
	service.AssertExpectations(t)
}
