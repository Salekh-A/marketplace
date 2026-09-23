package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"analytics/internal/model"

	"github.com/stretchr/testify/require"
)

func TestServer_GetRevenue(t *testing.T) {
	service := &mockAnalyticsService{
		getRevenueFunc: func(ctx context.Context) (model.Revenue, error) {
			return model.Revenue{
				Revenue: 1500,
			}, nil
		},
	}

	server := NewServer(service)

	response, err := server.GetRevenue(
		context.Background(),
		&GetRevenueRequest{},
	)

	require.NoError(t, err)
	require.Equal(t, 1500.0, response.Revenue)
}

func TestServer_GetRevenue_Error(t *testing.T) {
	expectedErr := errors.New("database error")

	service := &mockAnalyticsService{
		getRevenueFunc: func(ctx context.Context) (model.Revenue, error) {
			return model.Revenue{}, expectedErr
		},
	}

	server := NewServer(service)

	_, err := server.GetRevenue(
		context.Background(),
		&GetRevenueRequest{},
	)

	require.ErrorIs(t, err, expectedErr)
}

func TestServer_GetOrders(t *testing.T) {
	service := &mockAnalyticsService{
		getOrdersFunc: func(ctx context.Context) (model.OrderStats, error) {
			return model.OrderStats{
				Orders: 25,
			}, nil
		},
	}

	server := NewServer(service)

	response, err := server.GetOrders(
		context.Background(),
		&GetOrdersRequest{},
	)

	require.NoError(t, err)
	require.Equal(t, int32(25), response.Orders)
}

func TestServer_GetAverageCheck(t *testing.T) {
	service := &mockAnalyticsService{
		getAverageCheckFunc: func(ctx context.Context) (model.AverageCheck, error) {
			return model.AverageCheck{
				Average: 250.50,
			}, nil
		},
	}

	server := NewServer(service)

	response, err := server.GetAverageCheck(
		context.Background(),
		&GetAverageCheckRequest{},
	)

	require.NoError(t, err)
	require.Equal(t, 250.50, response.Average)
}

func TestServer_RecordOrder(t *testing.T) {
	called := false

	service := &mockAnalyticsService{
		recordOrderFunc: func(
			ctx context.Context,
			id string,
			productID string,
			quantity int,
			price float64,
			status string,
			createdAt time.Time,
		) error {
			called = true

			require.Equal(
				t,
				"11111111-1111-1111-1111-111111111111",
				id,
			)

			require.Equal(
				t,
				"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				productID,
			)

			require.Equal(t, 2, quantity)
			require.Equal(t, 100.0, price)
			require.Equal(t, "created", status)

			return nil
		},
	}

	server := NewServer(service)

	_, err := server.RecordOrder(
		context.Background(),
		&RecordOrderRequest{
			Id:        "11111111-1111-1111-1111-111111111111",
			ProductId: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			Quantity:  2,
			Price:     100,
			Status:    "created",
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	)

	require.NoError(t, err)
	require.True(t, called)
}

func TestServer_RecordOrder_Error(t *testing.T) {
	expectedErr := errors.New("database error")

	service := &mockAnalyticsService{
		recordOrderFunc: func(
			ctx context.Context,
			id string,
			productID string,
			quantity int,
			price float64,
			status string,
			createdAt time.Time,
		) error {
			return expectedErr
		},
	}

	server := NewServer(service)

	_, err := server.RecordOrder(
		context.Background(),
		&RecordOrderRequest{
			Id:        "11111111-1111-1111-1111-111111111111",
			ProductId: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			Quantity:  1,
			Price:     100,
			Status:    "created",
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	)

	require.ErrorIs(t, err, expectedErr)
}
