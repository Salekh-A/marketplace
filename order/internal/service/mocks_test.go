package service

import (
	"context"
	productpb "order/internal/grpc/product"
	"order/internal/model"
)

type mockOrderRepository struct {
	createFunc       func(context.Context, model.Order) error
	getByIDFunc      func(context.Context, string) (model.Order, error)
	updateStatusFunc func(context.Context, string, string) error
	deleteFunc       func(context.Context, string) error
}

func (m *mockOrderRepository) Create(
	ctx context.Context,
	order model.Order,
) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, order)
	}

	return nil
}

func (m *mockOrderRepository) GetByID(
	ctx context.Context,
	id string,
) (model.Order, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}

	return model.Order{}, nil
}

func (m *mockOrderRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, status)
	}

	return nil
}

func (m *mockOrderRepository) Delete(
	ctx context.Context,
	id string,
) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}

	return nil
}

type mockProductClient struct {
	getProductFunc    func(context.Context, string) (*productpb.GetProductResponse, error)
	decreaseStockFunc func(context.Context, string, int) error
}

func (m *mockProductClient) GetProduct(
	ctx context.Context,
	id string,
) (*productpb.GetProductResponse, error) {
	if m.getProductFunc != nil {
		return m.getProductFunc(ctx, id)
	}

	return &productpb.GetProductResponse{}, nil
}

func (m *mockProductClient) DecreaseStock(
	ctx context.Context,
	id string,
	quantity int,
) error {
	if m.decreaseStockFunc != nil {
		return m.decreaseStockFunc(ctx, id, quantity)
	}

	return nil
}

type mockAnalyticsClient struct {
	recordOrderFunc func(
		context.Context,
		string,
		string,
		int,
		float64,
		string,
		string,
	) error
}

func (m *mockAnalyticsClient) RecordOrder(
	ctx context.Context,
	id string,
	productID string,
	quantity int,
	price float64,
	status string,
	createdAt string,
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
