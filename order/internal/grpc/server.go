package grpc

import (
	"context"

	orderpb "marketplace-api/gen/order"
	"order/internal/model"
	"order/internal/service"
)

type OrderService interface {
	GetByID(
		ctx context.Context,
		id string,
	) (model.Order, error)

	Create(
		ctx context.Context,
		productID string,
		quantity int,
		price float64,
	) (model.Order, error)

	Cancel(
		ctx context.Context,
		id string,
	) error
}

type Server struct {
	orderpb.UnimplementedOrderServiceServer
	service OrderService
}

func NewServer(service OrderService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) GetOrder(
	ctx context.Context,
	req *orderpb.GetOrderRequest,
) (*orderpb.GetOrderResponse, error) {
	order, err := s.service.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &orderpb.GetOrderResponse{
		Id:        order.ID,
		ProductId: order.ProductID,
		Quantity:  int32(order.Quantity),
		Price:     order.Price,
		Status:    order.Status,
	}, nil
}

func (s *Server) CreateOrder(
	ctx context.Context,
	req *orderpb.CreateOrderRequest,
) (*orderpb.CreateOrderResponse, error) {
	order, err := s.service.Create(
		ctx,
		req.ProductId,
		int(req.Quantity),
		req.Price,
	)
	if err != nil {
		return nil, err
	}

	return &orderpb.CreateOrderResponse{
		Id:        order.ID,
		ProductId: order.ProductID,
		Quantity:  int32(order.Quantity),
		Price:     order.Price,
		Status:    order.Status,
	}, nil
}

func (s *Server) CancelOrder(
	ctx context.Context,
	req *orderpb.CancelOrderRequest,
) (*orderpb.CancelOrderResponse, error) {
	err := s.service.Cancel(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &orderpb.CancelOrderResponse{
		Success: true,
	}, nil
}

// compile-time check.
var _ OrderService = (*service.OrderService)(nil)
