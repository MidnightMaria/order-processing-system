package product

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

func (r *Repository) GetByID(ctx context.Context, id int64) (*Product, error) {
	product := &Product{}

	err := r.conn.QueryRow(
		ctx,
		`SELECT id, name, stock
		 FROM products
		 WHERE id = $1`,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Stock,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}