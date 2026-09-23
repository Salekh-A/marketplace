package service

import (
	"context"

	productpb "order/internal/grpc/product"
	"order/internal/model"
)

type OrderRepository interface {
	Create(
		ctx context.Context,
		order model.Order,
	) error

	GetByID(
		ctx context.Context,
		id string,
	) (model.Order, error)

	UpdateStatus(
		ctx context.Context,
		id string,
		status string,
	) error

	Delete(
		ctx context.Context,
		id string,
	) error
}

type ProductClient interface {
	GetProduct(
		ctx context.Context,
		id string,
	) (*productpb.GetProductResponse, error)

	DecreaseStock(
		ctx context.Context,
		id string,
		quantity int,
	) error
}

type AnalyticsClient interface {
	RecordOrder(
		ctx context.Context,
		id string,
		productID string,
		quantity int,
		price float64,
		status string,
		createdAt string,
	) error
}
