package order

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) CreateOrder(
	ctx context.Context,
	order *Order,
) error {
	return r.conn.QueryRow(
		ctx,
		`INSERT INTO orders (user_id, status, total_amount)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		order.UserID,
		order.Status,
		order.TotalAmount,
	).Scan(&order.ID)
}

func (r *Repository) CreateOrderTx(
	ctx context.Context,
	tx pgx.Tx,
	order *Order,
) error {
	return tx.QueryRow(
		ctx,
		`INSERT INTO orders (user_id, status, total_amount)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		order.UserID,
		order.Status,
		order.TotalAmount,
	).Scan(&order.ID)
}

