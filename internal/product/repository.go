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

func (r *Repository) GetByIDForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	id int64,
) (*Product, error) {
	product := &Product{}

	err := tx.QueryRow(
		ctx,
		`SELECT id, name, stock
		 FROM products
		 WHERE id = $1
		 FOR UPDATE`,
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

func (r *Repository) UpdateStockTx(
	ctx context.Context,
	tx pgx.Tx,
	id int64,
	stock int,
) error {
	_, err := tx.Exec(
		ctx,
		`UPDATE products
		 SET stock = $1
		 WHERE id = $2`,
		stock,
		id,
	)

	return err
}