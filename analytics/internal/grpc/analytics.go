package grpc

import (
	"context"
	"time"

	"analytics/internal/model"
	"analytics/internal/service"
)

type AnalyticsService interface {
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

type Server struct {
	UnimplementedAnalyticsServiceServer
	service AnalyticsService
}

func NewServer(service AnalyticsService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) GetRevenue(
	ctx context.Context,
	req *GetRevenueRequest,
) (*GetRevenueResponse, error) {
	result, err := s.service.GetRevenue(ctx)
	if err != nil {
		return nil, err
	}

	return &GetRevenueResponse{
		Revenue: result.Revenue,
	}, nil
}

func (s *Server) GetOrders(
	ctx context.Context,
	req *GetOrdersRequest,
) (*GetOrdersResponse, error) {
	result, err := s.service.GetOrders(ctx)
	if err != nil {
		return nil, err
	}

	return &GetOrdersResponse{
		Orders: int32(result.Orders),
	}, nil
}

func (s *Server) GetAverageCheck(
	ctx context.Context,
	req *GetAverageCheckRequest,
) (*GetAverageCheckResponse, error) {
	result, err := s.service.GetAverageCheck(ctx)
	if err != nil {
		return nil, err
	}

	return &GetAverageCheckResponse{
		Average: result.Average,
	}, nil
}

func (s *Server) RecordOrder(
	ctx context.Context,
	req *RecordOrderRequest,
) (*RecordOrderResponse, error) {
	createdAt, err := time.Parse(time.RFC3339, req.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = s.service.RecordOrder(
		ctx,
		req.Id,
		req.ProductId,
		int(req.Quantity),
		req.Price,
		req.Status,
		createdAt,
	)
	if err != nil {
		return nil, err
	}

	return &RecordOrderResponse{
		Success: true,
	}, nil
}

var _ AnalyticsService = (*service.AnalyticsService)(nil)
