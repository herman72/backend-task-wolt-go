package models

// VenueStaticResponse represents the static information of a venue,
// including its geographical location.
type VenueStaticResponse struct {
	VenueRaw *struct {
		Location struct {
			Coordinates []float64 `json:"coordinates"` // Coordinates of the venue [longitude, latitude].
		} `json:"location"`
	} `json:"venue_raw"`
}

// VenueDynamicResponse represents the dynamic information of a venue,
// including delivery specifications and pricing details.
type VenueDynamicResponse struct {
	VenueRaw *struct {
		DeliverySpecs struct {
			OrderMinimumNoSurcharge int `json:"order_minimum_no_surcharge"` // Minimum order value to avoid surcharge.
			DeliveryPricing         struct {
				BasePrice      int `json:"base_price"` // Base delivery price.
				DistanceRanges []struct {
					Min int     `json:"min"` // Minimum distance range (inclusive).
					Max int     `json:"max"` // Maximum distance range (exclusive).
					A   int     `json:"a"`   // Constant factor for delivery fee calculation.
					B   float64 `json:"b"`   // Multiplier factor for delivery fee calculation.
				} `json:"distance_ranges"` // List of distance ranges for pricing.
			} `json:"delivery_pricing"`
		} `json:"delivery_specs"`
	} `json:"venue_raw"`
}

// PriceResponse represents the response for delivery pricing calculations.
type PriceResponse struct {
	TotalPrice          int `json:"total_price"`          // Total price including cart value, delivery fee, and surcharge.
	SmallOrderSurcharge int `json:"small_order_surcharge"` // Surcharge for orders below the minimum value.
	CartValue           int `json:"cart_value"`           // Value of the cart.
	Delivery            struct {
		Fee      int `json:"fee"`      // Calculated delivery fee.
		Distance int `json:"distance"` // Distance between venue and user in meters.
	} `json:"delivery"`
}

// OrderInfo represents the information about an order required for delivery fee calculations.
type OrderInfo struct {
	Slug      string  `json:"slug"`      // Unique identifier for the venue.
	Lat       float64 `json:"lat"`       // Latitude of the user's location.
	Lon       float64 `json:"lon"`       // Longitude of the user's location.
	CartValue int     `json:"cart_value"` // Value of the user's cart.
}

// DistanceRange represents a range of distances and associated pricing factors.
type DistanceRange struct {
	Min int     `json:"min"` // Minimum distance range (inclusive).
	Max int     `json:"max"` // Maximum distance range (exclusive).
	A   int     `json:"a"`   // Constant factor for pricing.
	B   float64 `json:"b"`   // Multiplier factor for pricing.
}

// ServerError represents an error message to be sent to the client.
type ServerError struct {
	Msg string `json:"message"` // Error message.
}

// Config represents the configuration settings for the server and API.
type Config struct {
	Server struct {
		Port int `yaml:"port"` // Port number for the server to listen on.
	} `yaml:"server"`

	API struct {
		BaseURL string `yaml:"base_url"` // Base URL for external API calls.
	} `yaml:"api"`
}
