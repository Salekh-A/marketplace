package grpc

import (
	"context"

	"product/internal/model"
)

type mockProductService struct {
	getByIDFunc       func(context.Context, string) (model.Product, error)
	getAllFunc        func(context.Context) ([]model.Product, error)
	decreaseStockFunc func(context.Context, string, int) error
}

func (m *mockProductService) GetByID(
	ctx context.Context,
	id string,
) (model.Product, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}

	return model.Product{}, nil
}

func (m *mockProductService) GetAll(
	ctx context.Context,
) ([]model.Product, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}

	return nil, nil
}

func (m *mockProductService) DecreaseStock(
	ctx context.Context,
	id string,
	quantity int,
) error {
	if m.decreaseStockFunc != nil {
		return m.decreaseStockFunc(ctx, id, quantity)
	}

	return nil
}
