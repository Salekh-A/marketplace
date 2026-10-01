package grpc

import (
	"context"
	"time"

	"analytics/internal/model"
	"analytics/internal/service"
	analyticspb "marketplace-api/gen/analytics"
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
	analyticspb.UnimplementedAnalyticsServiceServer
	service AnalyticsService
}

func NewServer(service AnalyticsService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) GetRevenue(
	ctx context.Context,
	req *analyticspb.GetRevenueRequest,
) (*analyticspb.GetRevenueResponse, error) {
	result, err := s.service.GetRevenue(ctx)
	if err != nil {
		return nil, err
	}

	return &analyticspb.GetRevenueResponse{
		Revenue: result.Revenue,
	}, nil
}

func (s *Server) GetOrders(
	ctx context.Context,
	req *analyticspb.GetOrdersRequest,
) (*analyticspb.GetOrdersResponse, error) {
	result, err := s.service.GetOrders(ctx)
	if err != nil {
		return nil, err
	}

	return &analyticspb.GetOrdersResponse{
		Orders: int32(result.Orders),
	}, nil
}

func (s *Server) GetAverageCheck(
	ctx context.Context,
	req *analyticspb.GetAverageCheckRequest,
) (*analyticspb.GetAverageCheckResponse, error) {
	result, err := s.service.GetAverageCheck(ctx)
	if err != nil {
		return nil, err
	}

	return &analyticspb.GetAverageCheckResponse{
		Average: result.Average,
	}, nil
}

func (s *Server) RecordOrder(
	ctx context.Context,
	req *analyticspb.RecordOrderRequest,
) (*analyticspb.RecordOrderResponse, error) {
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

	return &analyticspb.RecordOrderResponse{
		Success: true,
	}, nil
}

var _ AnalyticsService = (*service.AnalyticsService)(nil)
