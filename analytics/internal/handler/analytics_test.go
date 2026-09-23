package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"analytics/internal/model"

	"github.com/stretchr/testify/require"
)

func TestAnalyticsHandler_Revenue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := &mockAnalyticsService{
			getRevenueFunc: func(ctx context.Context) (model.Revenue, error) {
				return model.Revenue{
					Revenue: 1500.50,
				}, nil
			},
		}

		handler := NewAnalyticsHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/analytics/revenue", nil)
		rec := httptest.NewRecorder()

		handler.Revenue(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "1500.5")
		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	})

	t.Run("service error", func(t *testing.T) {
		expectedErr := errors.New("database error")

		mockService := &mockAnalyticsService{
			getRevenueFunc: func(ctx context.Context) (model.Revenue, error) {
				return model.Revenue{}, expectedErr
			},
		}

		handler := NewAnalyticsHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/analytics/revenue", nil)
		rec := httptest.NewRecorder()

		handler.Revenue(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "failed to get revenue")
	})
}

func TestAnalyticsHandler_Orders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := &mockAnalyticsService{
			getOrdersFunc: func(ctx context.Context) (model.OrderStats, error) {
				return model.OrderStats{
					Orders: 42,
				}, nil
			},
		}

		handler := NewAnalyticsHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/analytics/orders", nil)
		rec := httptest.NewRecorder()

		handler.Orders(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "42")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockAnalyticsService{
			getOrdersFunc: func(ctx context.Context) (model.OrderStats, error) {
				return model.OrderStats{}, errors.New("database error")
			},
		}

		handler := NewAnalyticsHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/analytics/orders", nil)
		rec := httptest.NewRecorder()

		handler.Orders(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "failed to get orders")
	})
}

func TestAnalyticsHandler_AverageCheck(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := &mockAnalyticsService{
			getAverageCheckFunc: func(ctx context.Context) (model.AverageCheck, error) {
				return model.AverageCheck{
					Average: 375.25,
				}, nil
			},
		}

		handler := NewAnalyticsHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/analytics/average-check", nil)
		rec := httptest.NewRecorder()

		handler.AverageCheck(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "375.25")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockAnalyticsService{
			getAverageCheckFunc: func(ctx context.Context) (model.AverageCheck, error) {
				return model.AverageCheck{}, errors.New("database error")
			},
		}

		handler := NewAnalyticsHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/analytics/average-check", nil)
		rec := httptest.NewRecorder()

		handler.AverageCheck(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "failed to get average check")
	})
}

func TestAnalyticsHandler_TopProducts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := &mockAnalyticsService{
			getTopProductsFunc: func(ctx context.Context) ([]model.ProductStats, error) {
				return []model.ProductStats{
					{
						ProductID: "11111111-1111-1111-1111-111111111111",
						Revenue:   1000,
						Orders:    10,
					},
					{
						ProductID: "22222222-2222-2222-2222-222222222222",
						Revenue:   500,
						Orders:    5,
					},
				}, nil
			},
		}

		handler := NewAnalyticsHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/analytics/top-products", nil)
		rec := httptest.NewRecorder()

		handler.TopProducts(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "11111111-1111-1111-1111-111111111111")
		require.Contains(t, rec.Body.String(), "1000")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockAnalyticsService{
			getTopProductsFunc: func(ctx context.Context) ([]model.ProductStats, error) {
				return nil, errors.New("database error")
			},
		}

		handler := NewAnalyticsHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/analytics/top-products", nil)
		rec := httptest.NewRecorder()

		handler.TopProducts(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), "failed to get top products")
	})
}
