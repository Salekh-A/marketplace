package handler

import (
	"context"

	"analytics/internal/model"
)

type mockAnalyticsService struct {
	getRevenueFunc      func(context.Context) (model.Revenue, error)
	getOrdersFunc       func(context.Context) (model.OrderStats, error)
	getAverageCheckFunc func(context.Context) (model.AverageCheck, error)
	getTopProductsFunc  func(context.Context) ([]model.ProductStats, error)
}

func (m *mockAnalyticsService) GetRevenue(
	ctx context.Context,
) (model.Revenue, error) {
	return m.getRevenueFunc(ctx)
}

func (m *mockAnalyticsService) GetOrders(
	ctx context.Context,
) (model.OrderStats, error) {
	return m.getOrdersFunc(ctx)
}

func (m *mockAnalyticsService) GetAverageCheck(
	ctx context.Context,
) (model.AverageCheck, error) {
	return m.getAverageCheckFunc(ctx)
}

func (m *mockAnalyticsService) GetTopProducts(
	ctx context.Context,
) ([]model.ProductStats, error) {
	return m.getTopProductsFunc(ctx)
}

var _ AnalyticsService = (*mockAnalyticsService)(nil)
