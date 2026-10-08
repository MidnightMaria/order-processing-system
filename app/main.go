package main

import (
	"context"
	"fmt"

	"order-processing-system/internal/database"
	"order-processing-system/internal/product"
	"order-processing-system/internal/order"
)

func main() {
	ctx := context.Background()

	conn, err := database.Connect(ctx)
	if err != nil {
		panic(err)
	}

	defer conn.Close(ctx)

	productRepo := product.NewRepository(conn)
	productService := product.NewService(productRepo)

	orderRepo := order.NewRepository(conn)
	orderService := order.NewService(orderRepo, conn)

	newOrder := &order.Order{
	UserID:      1,
	Status:      "pending",
	TotalAmount: 50000,
	}

	err = orderService.CreateOrder(ctx, newOrder)
	if err != nil {
		panic(err)
	}

	fmt.Println("Order ID:", newOrder.ID)

	result, err := productService.GetByID(ctx, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println("ID:", result.ID)
	fmt.Println("Name:", result.Name)
	fmt.Println("Stock:", result.Stock)
}

