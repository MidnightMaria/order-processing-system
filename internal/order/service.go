package order

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	repository *Repository
	conn       *pgx.Conn
}

func NewService(repository *Repository, conn *pgx.Conn) *Service {
	return &Service{
		repository: repository,
		conn:       conn,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *Order) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	err = s.repository.CreateOrderTx(ctx, tx, order)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}