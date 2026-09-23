package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"analytics/internal/model"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func newTestService(
	t *testing.T,
	repository AnalyticsRepository,
) (*AnalyticsService, *miniredis.Miniredis) {
	t.Helper()

	miniRedis := miniredis.RunT(t)

	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	service := NewAnalyticsService(
		repository,
		redisClient,
	)

	return service, miniRedis
}

func TestAnalyticsService_GetRevenue_CacheMiss(t *testing.T) {
	repository := &mockAnalyticsRepository{
		getRevenueFunc: func(ctx context.Context) (model.Revenue, error) {
			return model.Revenue{
				Revenue: 1500,
			}, nil
		},
	}

	service, miniRedis := newTestService(t, repository)
	defer miniRedis.Close()

	result, err := service.GetRevenue(context.Background())

	require.NoError(t, err)
	require.Equal(t, 1500.0, result.Revenue)

	cached, err := miniRedis.Get("analytics:revenue")

	require.NoError(t, err)
	require.Contains(t, cached, "1500")
}

func TestAnalyticsService_GetRevenue_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	repository := &mockAnalyticsRepository{
		getRevenueFunc: func(ctx context.Context) (model.Revenue, error) {
			return model.Revenue{}, expectedErr
		},
	}

	service, miniRedis := newTestService(t, repository)
	defer miniRedis.Close()

	result, err := service.GetRevenue(context.Background())

	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, 0.0, result.Revenue)
}

func TestAnalyticsService_GetOrders_CacheMiss(t *testing.T) {
	repository := &mockAnalyticsRepository{
		getOrdersFunc: func(ctx context.Context) (model.OrderStats, error) {
			return model.OrderStats{
				Orders: 10,
			}, nil
		},
	}

	service, miniRedis := newTestService(t, repository)
	defer miniRedis.Close()

	result, err := service.GetOrders(context.Background())

	require.NoError(t, err)
	require.Equal(t, 10, result.Orders)

	cached, err := miniRedis.Get("analytics:orders")

	require.NoError(t, err)
	require.Contains(t, cached, "10")
}

func TestAnalyticsService_GetAverageCheck_CacheMiss(t *testing.T) {
	repository := &mockAnalyticsRepository{
		getAverageCheckFunc: func(ctx context.Context) (model.AverageCheck, error) {
			return model.AverageCheck{
				Average: 250.50,
			}, nil
		},
	}

	service, miniRedis := newTestService(t, repository)
	defer miniRedis.Close()

	result, err := service.GetAverageCheck(context.Background())

	require.NoError(t, err)
	require.Equal(t, 250.50, result.Average)

	cached, err := miniRedis.Get("analytics:average-check")

	require.NoError(t, err)
	require.Contains(t, cached, "250.5")
}

func TestAnalyticsService_GetTopProducts_CacheMiss(t *testing.T) {
	repository := &mockAnalyticsRepository{
		getTopProductsFunc: func(ctx context.Context) ([]model.ProductStats, error) {
			return []model.ProductStats{
				{
					ProductID: "product-1",
					Revenue:   1000,
					Orders:    5,
				},
			}, nil
		},
	}

	service, miniRedis := newTestService(t, repository)
	defer miniRedis.Close()

	result, err := service.GetTopProducts(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "product-1", result[0].ProductID)
	require.Equal(t, 1000.0, result[0].Revenue)
	require.Equal(t, 5, result[0].Orders)

	cached, err := miniRedis.Get("analytics:top-products")

	require.NoError(t, err)
	require.Contains(t, cached, "product-1")
}

func TestAnalyticsService_GetRevenue_CacheHit(t *testing.T) {
	repositoryCalled := false

	repository := &mockAnalyticsRepository{
		getRevenueFunc: func(ctx context.Context) (model.Revenue, error) {
			repositoryCalled = true

			return model.Revenue{
				Revenue: 9999,
			}, nil
		},
	}

	service, miniRedis := newTestService(t, repository)
	defer miniRedis.Close()

	err := miniRedis.Set(
		"analytics:revenue",
		`{"revenue":1500}`,
	)

	require.NoError(t, err)

	result, err := service.GetRevenue(context.Background())

	require.NoError(t, err)
	require.Equal(t, 1500.0, result.Revenue)
	require.False(t, repositoryCalled)
}

func TestAnalyticsService_RecordOrder(t *testing.T) {
	recordCalled := false

	repository := &mockAnalyticsRepository{
		recordOrderFunc: func(
			ctx context.Context,
			id string,
			productID string,
			quantity int,
			price float64,
			status string,
			createdAt time.Time,
		) error {
			recordCalled = true

			require.Equal(t, "order-1", id)
			require.Equal(t, "product-1", productID)
			require.Equal(t, 2, quantity)
			require.Equal(t, 100.0, price)
			require.Equal(t, "created", status)

			return nil
		},
	}

	service, miniRedis := newTestService(t, repository)
	defer miniRedis.Close()

	cacheKeys := []string{
		"analytics:revenue",
		"analytics:orders",
		"analytics:average-check",
		"analytics:top-products",
	}

	for _, key := range cacheKeys {
		err := miniRedis.Set(key, "cached-data")
		require.NoError(t, err)
	}

	err := service.RecordOrder(
		context.Background(),
		"order-1",
		"product-1",
		2,
		100.0,
		"created",
		time.Now(),
	)

	require.NoError(t, err)
	require.True(t, recordCalled)

	for _, key := range cacheKeys {
		require.False(t, miniRedis.Exists(key))
	}
}
