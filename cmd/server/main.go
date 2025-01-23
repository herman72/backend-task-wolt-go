package main

import (
	"backend-wolt-go/internal/api"
	"backend-wolt-go/internal/client"
	"backend-wolt-go/internal/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main(){
	venueProvider := service.NewVenueProvider("https://consumer-api.development.dev.woltapi.com/home-assignment-api/v1/venues")
	dopcService := client.NewDOPC(venueProvider)
	handler := api.NewHandler(dopcService)
	port := ":8000"
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v1/delivery-order-price", handler.GetDeliveryOrderPrice)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	log.Printf("Starting server on %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}
}