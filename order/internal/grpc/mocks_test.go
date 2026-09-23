package grpc

import (
	"context"

	"order/internal/model"
)

type mockOrderService struct {
	getByIDFunc func(
		context.Context,
		string,
	) (model.Order, error)

	createFunc func(
		context.Context,
		string,
		int,
		float64,
	) (model.Order, error)

	cancelFunc func(
		context.Context,
		string,
	) error
}

func (m *mockOrderService) GetByID(
	ctx context.Context,
	id string,
) (model.Order, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}

	return model.Order{}, nil
}

func (m *mockOrderService) Create(
	ctx context.Context,
	productID string,
	quantity int,
	price float64,
) (model.Order, error) {
	if m.createFunc != nil {
		return m.createFunc(
			ctx,
			productID,
			quantity,
			price,
		)
	}

	return model.Order{}, nil
}

func (m *mockOrderService) Cancel(
	ctx context.Context,
	id string,
) error {
	if m.cancelFunc != nil {
		return m.cancelFunc(ctx, id)
	}

	return nil
}
