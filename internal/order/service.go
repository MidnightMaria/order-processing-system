package order

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"order-processing-system/internal/product"
)

type Service struct {
	repository        *Repository
	productRepository *product.Repository
	conn              *pgx.Conn
}

func NewService(
	repository *Repository,
	productRepository *product.Repository,
	conn *pgx.Conn,
) *Service {
	return &Service{
		repository:        repository,
		productRepository: productRepository,
		conn:              conn,
	}
}

func (s *Service) CreateOrder(
	ctx context.Context,
	userID int64,
	productID int64,
	quantity int,
) (*Order, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	product, err := s.productRepository.GetByIDForUpdate(ctx, tx, productID)
	if err != nil {
		return nil, err
	}

	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}

	if product.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	newStock := product.Stock - quantity

	err = s.productRepository.UpdateStockTx(
		ctx,
		tx,
		productID,
		newStock,
	)
	if err != nil {
		return nil, err
	}
	// Hitung total harga order
	totalAmount := product.Price * float64(quantity)

	// Siapkan data order baru
	newOrder := &Order{
		UserID:      userID,
		Status:      "pending",
		TotalAmount: totalAmount,
	}

	err = s.repository.CreateOrderTx(ctx, tx, newOrder)
	if err != nil {
		return nil, err
	}

	newItem := &OrderItem{
		OrderID:   newOrder.ID,
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: product.Price,
	}

	err = s.repository.CreateOrderItemTx(ctx, tx, newItem)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return newOrder, nil
}
