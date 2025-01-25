package main

import (
	"backend-wolt-go/internal/api"
	"backend-wolt-go/internal/client"
	"backend-wolt-go/internal/service"
	"backend-wolt-go/internal/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// main is the entry point of the application. It initializes the configuration, services, and HTTP router,
// and starts the server.
func main() {
	// Load application configuration from the "config.yaml" file.
	config, err := utils.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize the venue provider service with the base URL from the configuration.
	venueProvider := service.NewVenueProvider(config.API.BaseURL)

	// Create a new DOPC (Delivery Order Price Calculator) service client using the venue provider.
	dopcService := client.NewDOPC(venueProvider)

	// Create a new API handler and pass the DOPC service to it.
	handler := api.NewHandler(dopcService)

	// Create a new router using the chi router package.
	r := chi.NewRouter()

	// Define an HTTP GET route for fetching delivery order prices.
	r.Get("/api/v1/delivery-order-price", handler.GetDeliveryOrderPrice)

	// Create an HTTP server instance with the specified address and handler.
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.Port),
		Handler: r,
	}

	// Log the server start message.
	log.Printf("Starting server on %d", config.Server.Port)

	// Start the HTTP server and handle any errors.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %d: %v\n", config.Server.Port, err)
	}
}
