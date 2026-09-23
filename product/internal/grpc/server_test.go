package grpc

import (
	"context"
	"errors"
	"testing"

	"product/internal/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_GetProduct_Success(t *testing.T) {
	expected := model.Product{
		ID:    "product-1",
		Name:  "iPhone",
		Price: 999.99,
		Stock: 10,
	}

	mockService := &mockProductService{
		getByIDFunc: func(
			ctx context.Context,
			id string,
		) (model.Product, error) {
			if id != "product-1" {
				t.Fatalf("unexpected product id: %s", id)
			}

			return expected, nil
		},
	}

	server := NewServer(mockService)

	response, err := server.GetProduct(
		context.Background(),
		&GetProductRequest{
			Id: "product-1",
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Id != expected.ID {
		t.Errorf("expected id %q, got %q", expected.ID, response.Id)
	}

	if response.Name != expected.Name {
		t.Errorf("expected name %q, got %q", expected.Name, response.Name)
	}

	if response.Price != expected.Price {
		t.Errorf("expected price %v, got %v", expected.Price, response.Price)
	}

	if response.Stock != int32(expected.Stock) {
		t.Errorf("expected stock %d, got %d", expected.Stock, response.Stock)
	}
}

func TestServer_GetProduct_Error(t *testing.T) {
	expectedErr := errors.New("product not found")

	mockService := &mockProductService{
		getByIDFunc: func(
			ctx context.Context,
			id string,
		) (model.Product, error) {
			return model.Product{}, expectedErr
		},
	}

	server := NewServer(mockService)

	response, err := server.GetProduct(
		context.Background(),
		&GetProductRequest{
			Id: "product-1",
		},
	)

	if response != nil {
		t.Fatal("expected nil response")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestServer_GetProducts_Success(t *testing.T) {
	expected := []model.Product{
		{
			ID:    "product-1",
			Name:  "iPhone",
			Price: 999.99,
			Stock: 10,
		},
		{
			ID:    "product-2",
			Name:  "MacBook",
			Price: 1999.99,
			Stock: 5,
		},
	}

	mockService := &mockProductService{
		getAllFunc: func(
			ctx context.Context,
		) ([]model.Product, error) {
			return expected, nil
		},
	}

	server := NewServer(mockService)

	response, err := server.GetProducts(
		context.Background(),
		&GetProductsRequest{},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(response.Products) != len(expected) {
		t.Fatalf(
			"expected %d products, got %d",
			len(expected),
			len(response.Products),
		)
	}

	for i, product := range expected {
		got := response.Products[i]

		if got.Id != product.ID {
			t.Errorf(
				"product %d: expected id %q, got %q",
				i,
				product.ID,
				got.Id,
			)
		}

		if got.Name != product.Name {
			t.Errorf(
				"product %d: expected name %q, got %q",
				i,
				product.Name,
				got.Name,
			)
		}

		if got.Price != product.Price {
			t.Errorf(
				"product %d: expected price %v, got %v",
				i,
				product.Price,
				got.Price,
			)
		}

		if got.Stock != int32(product.Stock) {
			t.Errorf(
				"product %d: expected stock %d, got %d",
				i,
				product.Stock,
				got.Stock,
			)
		}
	}
}

func TestServer_GetProducts_Error(t *testing.T) {
	expectedErr := errors.New("database error")

	mockService := &mockProductService{
		getAllFunc: func(
			ctx context.Context,
		) ([]model.Product, error) {
			return nil, expectedErr
		},
	}

	server := NewServer(mockService)

	response, err := server.GetProducts(
		context.Background(),
		&GetProductsRequest{},
	)

	if response != nil {
		t.Fatal("expected nil response")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestServer_DecreaseStock_Success(t *testing.T) {
	mockService := &mockProductService{
		decreaseStockFunc: func(
			ctx context.Context,
			id string,
			quantity int,
		) error {
			if id != "product-1" {
				t.Errorf("expected id product-1, got %s", id)
			}

			if quantity != 3 {
				t.Errorf("expected quantity 3, got %d", quantity)
			}

			return nil
		},
	}

	server := NewServer(mockService)

	response, err := server.DecreaseStock(
		context.Background(),
		&DecreaseStockRequest{
			Id:       "product-1",
			Quantity: 3,
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response == nil {
		t.Fatal("expected response")
	}

	if !response.Success {
		t.Fatal("expected success=true")
	}
}

func TestServer_DecreaseStock_Error(t *testing.T) {
	expectedErr := errors.New("not enough product stock")

	mockService := &mockProductService{
		decreaseStockFunc: func(
			ctx context.Context,
			id string,
			quantity int,
		) error {
			return expectedErr
		},
	}

	server := NewServer(mockService)

	response, err := server.DecreaseStock(
		context.Background(),
		&DecreaseStockRequest{
			Id:       "product-1",
			Quantity: 10,
		},
	)

	if response != nil {
		t.Fatal("expected nil response")
	}

	if err == nil {
		t.Fatal("expected error")
	}

	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf(
			"expected code %v, got %v",
			codes.InvalidArgument,
			status.Code(err),
		)
	}

	if status.Convert(err).Message() != expectedErr.Error() {
		t.Fatalf(
			"expected message %q, got %q",
			expectedErr.Error(),
			status.Convert(err).Message(),
		)
	}
}
