package service

import (
	"context"

	"product/internal/model"
)

type mockProductRepository struct {
	createFunc        func(context.Context, model.Product) error
	getByIDFunc       func(context.Context, string) (model.Product, error)
	getAllFunc        func(context.Context) ([]model.Product, error)
	updateFunc        func(context.Context, model.Product) error
	decreaseStockFunc func(context.Context, string, int) error
	deleteFunc        func(context.Context, string) error
}

func (m *mockProductRepository) Create(
	ctx context.Context,
	product model.Product,
) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, product)
	}

	return nil
}

func (m *mockProductRepository) GetByID(
	ctx context.Context,
	id string,
) (model.Product, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}

	return model.Product{}, nil
}

func (m *mockProductRepository) GetAll(
	ctx context.Context,
) ([]model.Product, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}

	return nil, nil
}

func (m *mockProductRepository) Update(
	ctx context.Context,
	product model.Product,
) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, product)
	}

	return nil
}

func (m *mockProductRepository) DecreaseStock(
	ctx context.Context,
	id string,
	quantity int,
) error {
	if m.decreaseStockFunc != nil {
		return m.decreaseStockFunc(ctx, id, quantity)
	}

	return nil
}

func (m *mockProductRepository) Delete(
	ctx context.Context,
	id string,
) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}

	return nil
}
