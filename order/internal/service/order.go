package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"order/internal/model"
)

type OrderService struct {
	repository      OrderRepository
	analyticsClient AnalyticsClient
	productClient   ProductClient
}

func NewOrderService(
	repository OrderRepository,
	analyticsClient AnalyticsClient,
	productClient ProductClient,
) *OrderService {
	return &OrderService{
		repository:      repository,
		analyticsClient: analyticsClient,
		productClient:   productClient,
	}
}

func (s *OrderService) Create(
	ctx context.Context,
	productID string,
	quantity int,
	price float64,
) (model.Order, error) {
	if productID == "" {
		return model.Order{}, errors.New("product_id is required")
	}

	if quantity <= 0 {
		return model.Order{}, errors.New("quantity must be greater than zero")
	}

	if price < 0 {
		return model.Order{}, errors.New("price cannot be negative")
	}

	_, err := s.productClient.GetProduct(ctx, productID)
	if err != nil {
		return model.Order{}, errors.New("product not found")
	}

	err = s.productClient.DecreaseStock(ctx, productID, quantity)
	if err != nil {
		return model.Order{}, err
	}

	now := time.Now()

	order := model.Order{
		ID:        uuid.New().String(),
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
		Status:    "created",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(ctx, order); err != nil {
		return model.Order{}, err
	}

	if err := s.analyticsClient.RecordOrder(
		ctx,
		order.ID,
		order.ProductID,
		order.Quantity,
		order.Price,
		order.Status,
		order.CreatedAt.Format(time.RFC3339),
	); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *OrderService) GetByID(
	ctx context.Context,
	id string,
) (model.Order, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *OrderService) Cancel(
	ctx context.Context,
	id string,
) error {
	return s.repository.UpdateStatus(ctx, id, "cancelled")
}

func (s *OrderService) Delete(
	ctx context.Context,
	id string,
) error {
	return s.repository.Delete(ctx, id)
}
