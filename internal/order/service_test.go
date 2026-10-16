package order

import (
	"context"
	"fmt"
	"testing"
	"time"

	"order-processing-system/internal/database"
	"order-processing-system/internal/product"
)

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()

	conn, err := database.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	productRepo := product.NewRepository(conn)
	orderRepo := NewRepository(conn)
	orderService := NewService(orderRepo, productRepo, conn)

	// Create a test user
	email := fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())

	var userID int64
	err = conn.QueryRow(ctx,
		`INSERT INTO users (name, email)
		 VALUES ($1, $2)
		 RETURNING id`,
		"Test User", email,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create a test product
	var productID int64
	err = conn.QueryRow(ctx,
		`INSERT INTO products (name, price, stock)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		"Test Product", 100.0, 10,
	).Scan(&productID)
	if err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	// Clean up test data after the test finishes
	defer func() {
		cleanupCtx := context.Background()

		_, err := conn.Exec(cleanupCtx,
			`DELETE FROM order_items WHERE product_id = $1`,
			productID,
		)
		if err != nil {
			t.Errorf("failed to clean up order items: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM orders WHERE user_id = $1`,
			userID,
		)
		if err != nil {
			t.Errorf("failed to clean up orders: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM products WHERE id = $1`,
			productID,
		)
		if err != nil {
			t.Errorf("failed to clean up product: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM users WHERE id = $1`,
			userID,
		)
		if err != nil {
			t.Errorf("failed to clean up user: %v", err)
		}
	}()

	// Create an order for 2 units
	createdOrder, err := orderService.CreateOrder(
		ctx,
		userID,
		productID,
		2,
	)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Check order ID
	if createdOrder.ID <= 0 {
		t.Errorf("expected a valid order ID, got %d", createdOrder.ID)
	}

	// Check order status
	if createdOrder.Status != "pending" {
		t.Errorf("expected status pending, got %s", createdOrder.Status)
	}

	// Check total amount
	if createdOrder.TotalAmount != 200.0 {
		t.Errorf("expected total amount 200, got %f", createdOrder.TotalAmount)
	}

		
	var orderID int64
	var savedProductID int64
	var quantity int
	var unitPrice float64

	err = conn.QueryRow(
		ctx,
		`SELECT order_id, product_id, quantity, unit_price
		FROM order_items
		WHERE order_id = $1`,
		createdOrder.ID,
	).Scan(
		&orderID,
		&savedProductID,
		&quantity,
		&unitPrice,
	)


	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	// Check saved order item
	if quantity != 2 {
		t.Errorf("expected quantity 2, got %d", quantity)
	}

	if orderID != createdOrder.ID {
		t.Errorf("expected order ID %d, got %d", createdOrder.ID, orderID)
	}

	if savedProductID != productID {
		t.Errorf("expected product ID %d, got %d", productID, savedProductID)
	}

	if unitPrice != 100.0 {
		t.Errorf("expected unit price 100, got %f", unitPrice)
	}

	// Check remaining product stock
	updatedProduct, err := productRepo.GetByID(ctx, productID)
	if err != nil {
		t.Fatalf("failed to get updated product: %v", err)
	}

	if updatedProduct.Stock != 8 {
		t.Errorf("expected stock 8, got %d", updatedProduct.Stock)
	}
}

func TestCreateOrderInsufficientStock(t *testing.T) {
	ctx := context.Background()

	conn, err := database.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	productRepo := product.NewRepository(conn)
	orderRepo := NewRepository(conn)
	orderService := NewService(orderRepo, productRepo, conn)

	// Create a test user
	email := fmt.Sprintf(
		"test-insufficient-%d@example.com",
		time.Now().UnixNano(),
	)

	var userID int64
	err = conn.QueryRow(ctx,
		`INSERT INTO users (name, email)
		 VALUES ($1, $2)
		 RETURNING id`,
		"Test User", email,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create a test product with stock 1
	var productID int64
	err = conn.QueryRow(ctx,
		`INSERT INTO products (name, price, stock)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		"Insufficient Stock Product", 100.0, 1,
	).Scan(&productID)
	if err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	// Clean up test data after the test finishes
	defer func() {
		cleanupCtx := context.Background()

		// Delete order items first, in case unexpected rows exist
		_, err := conn.Exec(cleanupCtx,
			`DELETE FROM order_items
			 WHERE order_id IN (
			     SELECT id FROM orders WHERE user_id = $1
			 )`,
			userID,
		)
		if err != nil {
			t.Errorf("failed to clean up order items: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM orders WHERE user_id = $1`,
			userID,
		)
		if err != nil {
			t.Errorf("failed to clean up orders: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM products WHERE id = $1`,
			productID,
		)
		if err != nil {
			t.Errorf("failed to clean up product: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM users WHERE id = $1`,
			userID,
		)
		if err != nil {
			t.Errorf("failed to clean up user: %v", err)
		}
	}()

	// Try to order 2 units when stock is only 1
	_, err = orderService.CreateOrder(
		ctx,
		userID,
		productID,
		2,
	)

	// An error is expected because stock is insufficient
	if err == nil {
		t.Fatal("expected an error because stock is insufficient")
	}

	// Check that stock remains unchanged
	updatedProduct, err := productRepo.GetByID(ctx, productID)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}

	if updatedProduct.Stock != 1 {
		t.Errorf("expected stock to remain 1, got %d", updatedProduct.Stock)
	}

	// Check that no order was created
	var orderCount int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM orders WHERE user_id = $1`,
		userID,
	).Scan(&orderCount)
	if err != nil {
		t.Fatalf("failed to count orders: %v", err)
	}

	if orderCount != 0 {
		t.Errorf("expected 0 orders, got %d", orderCount)
	}
}

func TestCreateOrderInvalidQuantity(t *testing.T) {
	ctx := context.Background()

	conn, err := database.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	productRepo := product.NewRepository(conn)
	orderRepo := NewRepository(conn)
	orderService := NewService(orderRepo, productRepo, conn)

	// Create a test user
	email := fmt.Sprintf(
		"test-invalid-%d@example.com",
		time.Now().UnixNano(),
	)

	var userID int64
	err = conn.QueryRow(ctx,
		`INSERT INTO users (name, email)
		 VALUES ($1, $2)
		 RETURNING id`,
		"Test User", email,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create a test product
	var productID int64
	err = conn.QueryRow(ctx,
		`INSERT INTO products (name, price, stock)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		"Invalid Quantity Product", 100.0, 10,
	).Scan(&productID)
	if err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	// Clean up test data
	defer func() {
		cleanupCtx := context.Background()

		_, err := conn.Exec(cleanupCtx,
			`DELETE FROM order_items
			 WHERE order_id IN (
			     SELECT id FROM orders WHERE user_id = $1
			 )`,
			userID,
		)
		if err != nil {
			t.Errorf("failed to clean up order items: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM orders WHERE user_id = $1`,
			userID,
		)
		if err != nil {
			t.Errorf("failed to clean up orders: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM products WHERE id = $1`,
			productID,
		)
		if err != nil {
			t.Errorf("failed to clean up product: %v", err)
		}

		_, err = conn.Exec(cleanupCtx,
			`DELETE FROM users WHERE id = $1`,
			userID,
		)
		if err != nil {
			t.Errorf("failed to clean up user: %v", err)
		}
	}()

	// Try to create an order with quantity 0
	_, err = orderService.CreateOrder(
		ctx,
		userID,
		productID,
		0,
	)

	if err == nil {
		t.Fatal("expected an error for quantity 0")
	}

	// Check that stock remains unchanged
	updatedProduct, err := productRepo.GetByID(ctx, productID)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}

	if updatedProduct.Stock != 10 {
		t.Errorf("expected stock to remain 10, got %d", updatedProduct.Stock)
	}

	// Check that no order was created
	var orderCount int
	err = conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM orders WHERE user_id = $1`,
		userID,
	).Scan(&orderCount)
	if err != nil {
		t.Fatalf("failed to count orders: %v", err)
	}

	if orderCount != 0 {
		t.Errorf("expected 0 orders, got %d", orderCount)
	}
}