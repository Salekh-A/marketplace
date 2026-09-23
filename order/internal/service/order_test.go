package service

import (
	"context"
	"errors"
	productpb "order/internal/grpc/product"
	"order/internal/model"
	"testing"
)

func TestOrderService_Create_Validation(t *testing.T) {
	service := &OrderService{}

	tests := []struct {
		name      string
		productID string
		quantity  int
		price     float64
		wantError string
	}{
		{
			name:      "empty product id",
			productID: "",
			quantity:  1,
			price:     100,
			wantError: "product_id is required",
		},
		{
			name:      "zero quantity",
			productID: "product-1",
			quantity:  0,
			price:     100,
			wantError: "quantity must be greater than zero",
		},
		{
			name:      "negative quantity",
			productID: "product-1",
			quantity:  -1,
			price:     100,
			wantError: "quantity must be greater than zero",
		},
		{
			name:      "negative price",
			productID: "product-1",
			quantity:  1,
			price:     -100,
			wantError: "price cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Create(
				context.Background(),
				tt.productID,
				tt.quantity,
				tt.price,
			)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != tt.wantError {
				t.Fatalf(
					"expected error %q, got %q",
					tt.wantError,
					err.Error(),
				)
			}
		})
	}
}

func TestOrderService_Create_Success(t *testing.T) {
	var savedOrder model.Order

	repository := &mockOrderRepository{
		createFunc: func(
			ctx context.Context,
			order model.Order,
		) error {
			savedOrder = order
			return nil
		},
	}

	productClient := &mockProductClient{
		getProductFunc: func(
			ctx context.Context,
			id string,
		) (*productpb.GetProductResponse, error) {
			return &productpb.GetProductResponse{
				Id:    id,
				Name:  "iPhone 15",
				Price: 1000,
				Stock: 5,
			}, nil
		},
	}

	analyticsClient := &mockAnalyticsClient{
		recordOrderFunc: func(
			ctx context.Context,
			id string,
			productID string,
			quantity int,
			price float64,
			status string,
			createdAt string,
		) error {
			return nil
		},
	}

	service := NewOrderService(
		repository,
		analyticsClient,
		productClient,
	)

	order, err := service.Create(
		context.Background(),
		"product-1",
		2,
		1000,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.ID == "" {
		t.Fatal("expected order ID")
	}

	if order.ProductID != "product-1" {
		t.Fatalf(
			"expected product ID %q, got %q",
			"product-1",
			order.ProductID,
		)
	}

	if order.Quantity != 2 {
		t.Fatalf(
			"expected quantity %d, got %d",
			2,
			order.Quantity,
		)
	}

	if order.Price != 1000 {
		t.Fatalf(
			"expected price %.2f, got %.2f",
			1000.0,
			order.Price,
		)
	}

	if order.Status != "created" {
		t.Fatalf(
			"expected status %q, got %q",
			"created",
			order.Status,
		)
	}

	if savedOrder.ID != order.ID {
		t.Fatal("expected repository to receive created order")
	}
}

func TestOrderService_Create_ProductNotFound(t *testing.T) {
	productClient := &mockProductClient{
		getProductFunc: func(
			ctx context.Context,
			id string,
		) (*productpb.GetProductResponse, error) {
			return nil, errors.New("product not found")
		},
	}

	repository := &mockOrderRepository{}
	analyticsClient := &mockAnalyticsClient{}

	service := NewOrderService(
		repository,
		analyticsClient,
		productClient,
	)

	_, err := service.Create(
		context.Background(),
		"product-1",
		2,
		1000,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "product not found" {
		t.Fatalf(
			"expected error %q, got %q",
			"product not found",
			err.Error(),
		)
	}
}

func TestOrderService_Create_DecreaseStockError(t *testing.T) {
	productClient := &mockProductClient{
		getProductFunc: func(
			ctx context.Context,
			id string,
		) (*productpb.GetProductResponse, error) {
			return &productpb.GetProductResponse{
				Id:    id,
				Name:  "iPhone 15",
				Price: 1000,
				Stock: 5,
			}, nil
		},
		decreaseStockFunc: func(
			ctx context.Context,
			id string,
			quantity int,
		) error {
			return errors.New("not enough product stock")
		},
	}

	repository := &mockOrderRepository{}
	analyticsClient := &mockAnalyticsClient{}

	service := NewOrderService(
		repository,
		analyticsClient,
		productClient,
	)

	_, err := service.Create(
		context.Background(),
		"product-1",
		10,
		1000,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "not enough product stock" {
		t.Fatalf(
			"expected error %q, got %q",
			"not enough product stock",
			err.Error(),
		)
	}
}

func TestOrderService_Create_RepositoryError(t *testing.T) {
	repository := &mockOrderRepository{
		createFunc: func(
			ctx context.Context,
			order model.Order,
		) error {
			return errors.New("database error")
		},
	}

	productClient := &mockProductClient{
		getProductFunc: func(
			ctx context.Context,
			id string,
		) (*productpb.GetProductResponse, error) {
			return &productpb.GetProductResponse{
				Id:    id,
				Name:  "iPhone 15",
				Price: 1000,
				Stock: 5,
			}, nil
		},
	}

	analyticsClient := &mockAnalyticsClient{}

	service := NewOrderService(
		repository,
		analyticsClient,
		productClient,
	)

	_, err := service.Create(
		context.Background(),
		"product-1",
		2,
		1000,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "database error" {
		t.Fatalf(
			"expected error %q, got %q",
			"database error",
			err.Error(),
		)
	}
}

func TestOrderService_Create_AnalyticsError(t *testing.T) {
	analyticsClient := &mockAnalyticsClient{
		recordOrderFunc: func(
			ctx context.Context,
			id string,
			productID string,
			quantity int,
			price float64,
			status string,
			createdAt string,
		) error {
			return errors.New("analytics error")
		},
	}

	productClient := &mockProductClient{
		getProductFunc: func(
			ctx context.Context,
			id string,
		) (*productpb.GetProductResponse, error) {
			return &productpb.GetProductResponse{
				Id:    id,
				Name:  "iPhone 15",
				Price: 1000,
				Stock: 5,
			}, nil
		},
	}

	repository := &mockOrderRepository{}

	service := NewOrderService(
		repository,
		analyticsClient,
		productClient,
	)

	_, err := service.Create(
		context.Background(),
		"product-1",
		2,
		1000,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "analytics error" {
		t.Fatalf(
			"expected error %q, got %q",
			"analytics error",
			err.Error(),
		)
	}
}

func TestOrderService_GetByID(t *testing.T) {
	expected := model.Order{
		ID:        "order-1",
		ProductID: "product-1",
		Quantity:  2,
		Price:     1000,
		Status:    "created",
	}

	repository := &mockOrderRepository{
		getByIDFunc: func(
			ctx context.Context,
			id string,
		) (model.Order, error) {
			return expected, nil
		},
	}

	service := NewOrderService(
		repository,
		&mockAnalyticsClient{},
		&mockProductClient{},
	)

	order, err := service.GetByID(
		context.Background(),
		"order-1",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.ID != expected.ID {
		t.Fatalf(
			"expected ID %q, got %q",
			expected.ID,
			order.ID,
		)
	}
}

func TestOrderService_Cancel(t *testing.T) {
	var gotID string
	var gotStatus string

	repository := &mockOrderRepository{
		updateStatusFunc: func(
			ctx context.Context,
			id string,
			status string,
		) error {
			gotID = id
			gotStatus = status
			return nil
		},
	}

	service := NewOrderService(
		repository,
		&mockAnalyticsClient{},
		&mockProductClient{},
	)

	err := service.Cancel(
		context.Background(),
		"order-1",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotID != "order-1" {
		t.Fatalf("expected ID %q, got %q", "order-1", gotID)
	}

	if gotStatus != "cancelled" {
		t.Fatalf(
			"expected status %q, got %q",
			"cancelled",
			gotStatus,
		)
	}
}

func TestOrderService_Delete(t *testing.T) {
	var gotID string

	repository := &mockOrderRepository{
		deleteFunc: func(
			ctx context.Context,
			id string,
		) error {
			gotID = id
			return nil
		},
	}

	service := NewOrderService(
		repository,
		&mockAnalyticsClient{},
		&mockProductClient{},
	)

	err := service.Delete(
		context.Background(),
		"order-1",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotID != "order-1" {
		t.Fatalf(
			"expected ID %q, got %q",
			"order-1",
			gotID,
		)
	}
}
