package main

import (
	"backend-wolt-go/internal/api"
	"log"
	"net/http"
)

func main(){
	http.HandleFunc("/api/v1/delivery-order-price", api.Handler)
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}