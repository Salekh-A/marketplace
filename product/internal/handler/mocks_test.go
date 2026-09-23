package handler

import (
	"context"

	"product/internal/model"
)

type mockProductService struct {
	createFunc  func(context.Context, string, float64, int) (model.Product, error)
	getAllFunc  func(context.Context) ([]model.Product, error)
	getByIDFunc func(context.Context, string) (model.Product, error)
	updateFunc  func(context.Context, string, string, float64, int) (model.Product, error)
	deleteFunc  func(context.Context, string) error
}

func (m *mockProductService) Create(
	ctx context.Context,
	name string,
	price float64,
	stock int,
) (model.Product, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, name, price, stock)
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

func (m *mockProductService) GetByID(
	ctx context.Context,
	id string,
) (model.Product, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}

	return model.Product{}, nil
}

func (m *mockProductService) Update(
	ctx context.Context,
	id string,
	name string,
	price float64,
	stock int,
) (model.Product, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, name, price, stock)
	}

	return model.Product{}, nil
}

func (m *mockProductService) Delete(
	ctx context.Context,
	id string,
) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}

	return nil
}
