package grpc

import (
	"context"
	"time"

	"analytics/internal/model"
)

type mockAnalyticsService struct {
	getRevenueFunc      func(context.Context) (model.Revenue, error)
	getOrdersFunc       func(context.Context) (model.OrderStats, error)
	getAverageCheckFunc func(context.Context) (model.AverageCheck, error)
	getTopProductsFunc  func(context.Context) ([]model.ProductStats, error)
	recordOrderFunc     func(
		context.Context,
		string,
		string,
		int,
		float64,
		string,
		time.Time,
	) error
}

func (m *mockAnalyticsService) GetRevenue(
	ctx context.Context,
) (model.Revenue, error) {
	if m.getRevenueFunc != nil {
		return m.getRevenueFunc(ctx)
	}

	return model.Revenue{}, nil
}

func (m *mockAnalyticsService) GetOrders(
	ctx context.Context,
) (model.OrderStats, error) {
	if m.getOrdersFunc != nil {
		return m.getOrdersFunc(ctx)
	}

	return model.OrderStats{}, nil
}

func (m *mockAnalyticsService) GetAverageCheck(
	ctx context.Context,
) (model.AverageCheck, error) {
	if m.getAverageCheckFunc != nil {
		return m.getAverageCheckFunc(ctx)
	}

	return model.AverageCheck{}, nil
}

func (m *mockAnalyticsService) GetTopProducts(
	ctx context.Context,
) ([]model.ProductStats, error) {
	if m.getTopProductsFunc != nil {
		return m.getTopProductsFunc(ctx)
	}

	return nil, nil
}

func (m *mockAnalyticsService) RecordOrder(
	ctx context.Context,
	id string,
	productID string,
	quantity int,
	price float64,
	status string,
	createdAt time.Time,
) error {
	if m.recordOrderFunc != nil {
		return m.recordOrderFunc(
			ctx,
			id,
			productID,
			quantity,
			price,
			status,
			createdAt,
		)
	}

	return nil
}
