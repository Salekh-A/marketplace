package service

import (
	"context"
	"time"

	"analytics/internal/model"

	"github.com/redis/go-redis/v9"
)

type AnalyticsRepository interface {
	GetRevenue(ctx context.Context) (model.Revenue, error)
	GetOrders(ctx context.Context) (model.OrderStats, error)
	GetAverageCheck(ctx context.Context) (model.AverageCheck, error)
	GetTopProducts(ctx context.Context) ([]model.ProductStats, error)

	RecordOrder(
		ctx context.Context,
		id string,
		productID string,
		quantity int,
		price float64,
		status string,
		createdAt time.Time,
	) error
}

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd

	Set(
		ctx context.Context,
		key string,
		value interface{},
		expiration time.Duration,
	) *redis.StatusCmd

	Del(
		ctx context.Context,
		keys ...string,
	) *redis.IntCmd
}
