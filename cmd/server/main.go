package main

import (
	"backend-wolt-go/internal/api"
	"backend-wolt-go/internal/client"
	"backend-wolt-go/internal/service"
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main(){
	port := ":8000"
	
	// Initialize dependencies
	apiClient := client.NewAPIClient()
	calculatorService := service.NewCalculatorService()
	handler := api.NewHandler(apiClient, calculatorService)

	// Set up router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v1/delivery-order-price", handler.GetDeliveryOrderPrice)

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}