package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	conn, err := pgx.Connect(
		ctx,
		"postgres://postgres:postgres@localhost:5432/order_processing",
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}