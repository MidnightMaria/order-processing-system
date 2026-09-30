package main

import (
	"context"
	"fmt"

	"order-processing-system/internal/database"
	"order-processing-system/internal/product"
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

	result, err := productService.GetByID(ctx, 1)
	if err != nil {
		panic(err)
	}

	fmt.Println("ID:", result.ID)
	fmt.Println("Name:", result.Name)
	fmt.Println("Stock:", result.Stock)
}

