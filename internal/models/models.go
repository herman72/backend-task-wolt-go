package models


type VenueStaticResponse struct {
	VenueRaw *struct {
		Location struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"location"`
	} `json:"venue_raw"`
}

type VenueDynamicResponse struct {
	VenueRaw *struct {
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

type PriceResponse struct {
	TotalPrice          int `json:"total_price"`
	SmallOrderSurcharge int `json:"small_order_surcharge"`
	CartValue           int `json:"cart_value"`
	Delivery            struct {
		Fee      int `json:"fee"`
		Distance int `json:"distance"`
	} `json:"delivery"`
}

type OrderInfo struct {
	Slug string `json:"slug"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	CartValue int `json:"cart_value"`
}

type DistanceRange struct {
	Min int     `json:"min"`
	Max int     `json:"max"`
	A   int     `json:"a"`
	B   float64 `json:"b"`
}

type ServerError struct {
	Msg string `json:"message"`
}

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	API struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"api"`

}
