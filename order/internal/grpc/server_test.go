package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"order/internal/model"
)

func TestOrderServer_GetOrder(t *testing.T) {
	order := model.Order{
		ID:        "order-1",
		ProductID: "product-1",
		Quantity:  2,
		Price:     999.99,
		Status:    "created",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		service := &mockOrderService{
			getByIDFunc: func(
				ctx context.Context,
				id string,
			) (model.Order, error) {
				if id != order.ID {
					t.Errorf("expected id %q, got %q", order.ID, id)
				}

				return order, nil
			},
		}

		server := NewServer(service)

		response, err := server.GetOrder(
			context.Background(),
			&GetOrderRequest{
				Id: order.ID,
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if response.Id != order.ID {
			t.Errorf("expected id %q, got %q", order.ID, response.Id)
		}

		if response.ProductId != order.ProductID {
			t.Errorf(
				"expected product_id %q, got %q",
				order.ProductID,
				response.ProductId,
			)
		}

		if response.Quantity != int32(order.Quantity) {
			t.Errorf(
				"expected quantity %d, got %d",
				order.Quantity,
				response.Quantity,
			)
		}

		if response.Price != order.Price {
			t.Errorf(
				"expected price %v, got %v",
				order.Price,
				response.Price,
			)
		}

		if response.Status != order.Status {
			t.Errorf(
				"expected status %q, got %q",
				order.Status,
				response.Status,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		expectedErr := errors.New("order not found")

		service := &mockOrderService{
			getByIDFunc: func(
				ctx context.Context,
				id string,
			) (model.Order, error) {
				return model.Order{}, expectedErr
			},
		}

		server := NewServer(service)

		response, err := server.GetOrder(
			context.Background(),
			&GetOrderRequest{
				Id: order.ID,
			},
		)

		if response != nil {
			t.Fatal("expected nil response")
		}

		if !errors.Is(err, expectedErr) {
			t.Fatalf(
				"expected error %v, got %v",
				expectedErr,
				err,
			)
		}
	})
}

func TestOrderServer_CreateOrder(t *testing.T) {
	order := model.Order{
		ID:        "order-1",
		ProductID: "product-1",
		Quantity:  3,
		Price:     1499.99,
		Status:    "created",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		service := &mockOrderService{
			createFunc: func(
				ctx context.Context,
				productID string,
				quantity int,
				price float64,
			) (model.Order, error) {
				if productID != order.ProductID {
					t.Errorf(
						"expected product_id %q, got %q",
						order.ProductID,
						productID,
					)
				}

				if quantity != order.Quantity {
					t.Errorf(
						"expected quantity %d, got %d",
						order.Quantity,
						quantity,
					)
				}

				if price != order.Price {
					t.Errorf(
						"expected price %v, got %v",
						order.Price,
						price,
					)
				}

				return order, nil
			},
		}

		server := NewServer(service)

		response, err := server.CreateOrder(
			context.Background(),
			&CreateOrderRequest{
				ProductId: order.ProductID,
				Quantity:  int32(order.Quantity),
				Price:     order.Price,
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if response.Id != order.ID {
			t.Errorf(
				"expected id %q, got %q",
				order.ID,
				response.Id,
			)
		}

		if response.ProductId != order.ProductID {
			t.Errorf(
				"expected product_id %q, got %q",
				order.ProductID,
				response.ProductId,
			)
		}

		if response.Quantity != int32(order.Quantity) {
			t.Errorf(
				"expected quantity %d, got %d",
				order.Quantity,
				response.Quantity,
			)
		}

		if response.Price != order.Price {
			t.Errorf(
				"expected price %v, got %v",
				order.Price,
				response.Price,
			)
		}

		if response.Status != order.Status {
			t.Errorf(
				"expected status %q, got %q",
				order.Status,
				response.Status,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		expectedErr := errors.New("failed to create order")

		service := &mockOrderService{
			createFunc: func(
				ctx context.Context,
				productID string,
				quantity int,
				price float64,
			) (model.Order, error) {
				return model.Order{}, expectedErr
			},
		}

		server := NewServer(service)

		response, err := server.CreateOrder(
			context.Background(),
			&CreateOrderRequest{
				ProductId: order.ProductID,
				Quantity:  int32(order.Quantity),
				Price:     order.Price,
			},
		)

		if response != nil {
			t.Fatal("expected nil response")
		}

		if !errors.Is(err, expectedErr) {
			t.Fatalf(
				"expected error %v, got %v",
				expectedErr,
				err,
			)
		}
	})
}

func TestOrderServer_CancelOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mockOrderService{
			cancelFunc: func(
				ctx context.Context,
				id string,
			) error {
				if id != "order-1" {
					t.Errorf(
						"expected id %q, got %q",
						"order-1",
						id,
					)
				}

				return nil
			},
		}

		server := NewServer(service)

		response, err := server.CancelOrder(
			context.Background(),
			&CancelOrderRequest{
				Id: "order-1",
			},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !response.Success {
			t.Error("expected success to be true")
		}
	})

	t.Run("service error", func(t *testing.T) {
		expectedErr := status.Error(
			codes.NotFound,
			"order not found",
		)

		service := &mockOrderService{
			cancelFunc: func(
				ctx context.Context,
				id string,
			) error {
				return expectedErr
			},
		}

		server := NewServer(service)

		response, err := server.CancelOrder(
			context.Background(),
			&CancelOrderRequest{
				Id: "order-1",
			},
		)

		if response != nil {
			t.Fatal("expected nil response")
		}

		if status.Code(err) != codes.NotFound {
			t.Fatalf(
				"expected code %v, got %v",
				codes.NotFound,
				status.Code(err),
			)
		}
	})
}
