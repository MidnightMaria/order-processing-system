package main

import (
	"context"
	"log"
	"net/http"

	"order-processing-system/internal/database"
	"order-processing-system/internal/httpapi"
	"order-processing-system/internal/order"
	"order-processing-system/internal/product"
)

func main() {
	ctx := context.Background()

	conn, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	productRepo := product.NewRepository(conn)
	orderRepo := order.NewRepository(conn)

	orderService := order.NewService(orderRepo, productRepo, conn)
	orderHandler := httpapi.NewOrderHandler(orderService)

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", orderHandler.CreateOrder)

	log.Println("Server running at http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
