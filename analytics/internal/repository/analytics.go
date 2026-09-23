package repository

import (
	"context"
	"time"

	"analytics/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalyticsRepository struct {
	db *pgxpool.Pool
}

func NewAnalyticsRepository(db *pgxpool.Pool) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

func (r *AnalyticsRepository) GetRevenue(
	ctx context.Context,
) (model.Revenue, error) {
	var result model.Revenue

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COALESCE(SUM(price * quantity), 0)
		FROM analytics_orders
		WHERE status != 'cancelled'
		`,
	).Scan(&result.Revenue)

	return result, err
}

func (r *AnalyticsRepository) GetOrders(
	ctx context.Context,
) (model.OrderStats, error) {
	var result model.OrderStats

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM analytics_orders
		WHERE status != 'cancelled'
		`,
	).Scan(&result.Orders)

	return result, err
}

func (r *AnalyticsRepository) GetAverageCheck(
	ctx context.Context,
) (model.AverageCheck, error) {
	var result model.AverageCheck

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COALESCE(AVG(total), 0)
		FROM (
			SELECT SUM(price * quantity) AS total
			FROM analytics_orders
			WHERE status != 'cancelled'
			GROUP BY id
		) orders
		`,
	).Scan(&result.Average)

	return result, err
}

func (r *AnalyticsRepository) GetTopProducts(
	ctx context.Context,
) ([]model.ProductStats, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			product_id,
			SUM(price * quantity) AS revenue,
			COUNT(*) AS analytics_orders
		FROM analytics_orders
		WHERE status != 'cancelled'
		GROUP BY product_id
		ORDER BY revenue DESC
		LIMIT 10
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.ProductStats

	for rows.Next() {
		var stats model.ProductStats

		if err := rows.Scan(
			&stats.ProductID,
			&stats.Revenue,
			&stats.Orders,
		); err != nil {
			return nil, err
		}

		result = append(result, stats)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AnalyticsRepository) RecordOrder(
	ctx context.Context,
	id string,
	productID string,
	quantity int,
	price float64,
	status string,
	createdAt time.Time,
) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO analytics_orders (
			id,
			product_id,
			quantity,
			price,
			status,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		id,
		productID,
		quantity,
		price,
		status,
		createdAt,
	)

	return err
}
