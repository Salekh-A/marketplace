package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"order/internal/model"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order model.Order) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO orders (
			id,
			product_id,
			quantity,
			price,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		order.ID,
		order.ProductID,
		order.Quantity,
		order.Price,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (model.Order, error) {
	var order model.Order

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			product_id,
			quantity,
			price,
			status,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		`,
		id,
	).Scan(
		&order.ID,
		&order.ProductID,
		&order.Quantity,
		&order.Price,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	return order, err
}

func (r *OrderRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE orders
		SET status = $1,
			updated_at = NOW()
		WHERE id = $2
		`,
		status,
		id,
	)

	return err
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(
		ctx,
		`DELETE FROM orders WHERE id = $1`,
		id,
	)

	return err
}
