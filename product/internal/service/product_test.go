package service

import (
	"context"
	"errors"
	"testing"

	"product/internal/model"
)

func TestProductService_Create_Validation(t *testing.T) {
	service := &ProductService{}

	tests := []struct {
		name        string
		productName string
		price       float64
		stock       int
		wantError   string
	}{
		{
			name:        "empty name",
			productName: "",
			price:       100,
			stock:       10,
			wantError:   "name is required",
		},
		{
			name:        "whitespace name",
			productName: "   ",
			price:       100,
			stock:       10,
			wantError:   "name is required",
		},
		{
			name:        "negative price",
			productName: "Phone",
			price:       -100,
			stock:       10,
			wantError:   "price cannot be negative",
		},
		{
			name:        "negative stock",
			productName: "Phone",
			price:       100,
			stock:       -1,
			wantError:   "stock cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Create(
				context.Background(),
				tt.productName,
				tt.price,
				tt.stock,
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

func TestProductService_Create_Success(t *testing.T) {
	var savedProduct model.Product

	repository := &mockProductRepository{
		createFunc: func(
			ctx context.Context,
			product model.Product,
		) error {
			savedProduct = product
			return nil
		},
	}

	service := NewProductService(repository)

	product, err := service.Create(
		context.Background(),
		"  iPhone  ",
		999.99,
		10,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.ID == "" {
		t.Fatal("expected product ID, got empty")
	}

	if product.Name != "iPhone" {
		t.Fatalf(
			"expected name %q, got %q",
			"iPhone",
			product.Name,
		)
	}

	if product.Price != 999.99 {
		t.Fatalf(
			"expected price %.2f, got %.2f",
			999.99,
			product.Price,
		)
	}

	if product.Stock != 10 {
		t.Fatalf(
			"expected stock %d, got %d",
			10,
			product.Stock,
		)
	}

	if savedProduct.ID != product.ID {
		t.Fatal("repository received different product ID")
	}

	if savedProduct.Name != product.Name {
		t.Fatal("repository received different product name")
	}
}

func TestProductService_Create_RepositoryError(t *testing.T) {
	repository := &mockProductRepository{
		createFunc: func(
			ctx context.Context,
			product model.Product,
		) error {
			return errors.New("database error")
		},
	}

	service := NewProductService(repository)

	_, err := service.Create(
		context.Background(),
		"Phone",
		100,
		10,
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

func TestProductService_GetByID(t *testing.T) {
	expectedProduct := model.Product{
		ID:    "product-1",
		Name:  "Phone",
		Price: 999.99,
		Stock: 10,
	}

	repository := &mockProductRepository{
		getByIDFunc: func(
			ctx context.Context,
			id string,
		) (model.Product, error) {
			if id != "product-1" {
				t.Fatalf(
					"expected id %q, got %q",
					"product-1",
					id,
				)
			}

			return expectedProduct, nil
		},
	}

	service := NewProductService(repository)

	product, err := service.GetByID(
		context.Background(),
		"product-1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if product.ID != expectedProduct.ID {
		t.Fatalf(
			"expected ID %q, got %q",
			expectedProduct.ID,
			product.ID,
		)
	}

	if product.Name != expectedProduct.Name {
		t.Fatalf(
			"expected name %q, got %q",
			expectedProduct.Name,
			product.Name,
		)
	}
}

func TestProductService_GetByID_Error(t *testing.T) {
	repository := &mockProductRepository{
		getByIDFunc: func(
			ctx context.Context,
			id string,
		) (model.Product, error) {
			return model.Product{}, errors.New("product not found")
		},
	}

	service := NewProductService(repository)

	_, err := service.GetByID(
		context.Background(),
		"product-1",
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

func TestProductService_GetAll(t *testing.T) {
	expectedProducts := []model.Product{
		{
			ID:    "product-1",
			Name:  "Phone",
			Price: 999.99,
			Stock: 10,
		},
		{
			ID:    "product-2",
			Name:  "Laptop",
			Price: 1999.99,
			Stock: 5,
		},
	}

	repository := &mockProductRepository{
		getAllFunc: func(
			ctx context.Context,
		) ([]model.Product, error) {
			return expectedProducts, nil
		},
	}

	service := NewProductService(repository)

	products, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf(
			"expected %d products, got %d",
			2,
			len(products),
		)
	}

	if products[0].ID != "product-1" {
		t.Fatalf(
			"expected first product %q, got %q",
			"product-1",
			products[0].ID,
		)
	}
}

func TestProductService_GetAll_Error(t *testing.T) {
	repository := &mockProductRepository{
		getAllFunc: func(
			ctx context.Context,
		) ([]model.Product, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewProductService(repository)

	_, err := service.GetAll(context.Background())
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

func TestProductService_Update_Validation(t *testing.T) {
	service := &ProductService{}

	tests := []struct {
		name        string
		productName string
		price       float64
		stock       int
		wantError   string
	}{
		{
			name:        "empty name",
			productName: "",
			price:       100,
			stock:       10,
			wantError:   "name is required",
		},
		{
			name:        "whitespace name",
			productName: "   ",
			price:       100,
			stock:       10,
			wantError:   "name is required",
		},
		{
			name:        "negative price",
			productName: "Phone",
			price:       -1,
			stock:       10,
			wantError:   "price cannot be negative",
		},
		{
			name:        "negative stock",
			productName: "Phone",
			price:       100,
			stock:       -1,
			wantError:   "stock cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Update(
				context.Background(),
				"product-1",
				tt.productName,
				tt.price,
				tt.stock,
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

func TestProductService_Update_Success(t *testing.T) {
	var updatedProduct model.Product

	repository := &mockProductRepository{
		updateFunc: func(
			ctx context.Context,
			product model.Product,
		) error {
			updatedProduct = product
			return nil
		},
		getByIDFunc: func(
			ctx context.Context,
			id string,
		) (model.Product, error) {
			return model.Product{
				ID:    id,
				Name:  "Updated Phone",
				Price: 1200,
				Stock: 20,
			}, nil
		},
	}

	service := NewProductService(repository)

	product, err := service.Update(
		context.Background(),
		"product-1",
		"  Updated Phone  ",
		1200,
		20,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updatedProduct.ID != "product-1" {
		t.Fatalf(
			"expected ID %q, got %q",
			"product-1",
			updatedProduct.ID,
		)
	}

	if updatedProduct.Name != "Updated Phone" {
		t.Fatalf(
			"expected name %q, got %q",
			"Updated Phone",
			updatedProduct.Name,
		)
	}

	if product.Name != "Updated Phone" {
		t.Fatalf(
			"expected returned name %q, got %q",
			"Updated Phone",
			product.Name,
		)
	}

	if product.Price != 1200 {
		t.Fatalf(
			"expected price %.2f, got %.2f",
			1200.0,
			product.Price,
		)
	}

	if product.Stock != 20 {
		t.Fatalf(
			"expected stock %d, got %d",
			20,
			product.Stock,
		)
	}
}

func TestProductService_Update_RepositoryError(t *testing.T) {
	repository := &mockProductRepository{
		updateFunc: func(
			ctx context.Context,
			product model.Product,
		) error {
			return errors.New("database error")
		},
	}

	service := NewProductService(repository)

	_, err := service.Update(
		context.Background(),
		"product-1",
		"Phone",
		100,
		10,
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

func TestProductService_Update_GetByIDError(t *testing.T) {
	repository := &mockProductRepository{
		updateFunc: func(
			ctx context.Context,
			product model.Product,
		) error {
			return nil
		},
		getByIDFunc: func(
			ctx context.Context,
			id string,
		) (model.Product, error) {
			return model.Product{}, errors.New("product not found")
		},
	}

	service := NewProductService(repository)

	_, err := service.Update(
		context.Background(),
		"product-1",
		"Phone",
		100,
		10,
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

func TestProductService_DecreaseStock_Validation(t *testing.T) {
	service := &ProductService{}

	err := service.DecreaseStock(
		context.Background(),
		"product-1",
		0,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "quantity must be greater than zero" {
		t.Fatalf(
			"expected error %q, got %q",
			"quantity must be greater than zero",
			err.Error(),
		)
	}
}

func TestProductService_DecreaseStock_Success(t *testing.T) {
	var savedID string
	var savedQuantity int

	repository := &mockProductRepository{
		decreaseStockFunc: func(
			ctx context.Context,
			id string,
			quantity int,
		) error {
			savedID = id
			savedQuantity = quantity
			return nil
		},
	}

	service := NewProductService(repository)

	err := service.DecreaseStock(
		context.Background(),
		"product-1",
		3,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if savedID != "product-1" {
		t.Fatalf(
			"expected ID %q, got %q",
			"product-1",
			savedID,
		)
	}

	if savedQuantity != 3 {
		t.Fatalf(
			"expected quantity %d, got %d",
			3,
			savedQuantity,
		)
	}
}

func TestProductService_DecreaseStock_Error(t *testing.T) {
	repository := &mockProductRepository{
		decreaseStockFunc: func(
			ctx context.Context,
			id string,
			quantity int,
		) error {
			return errors.New("not enough product stock")
		},
	}

	service := NewProductService(repository)

	err := service.DecreaseStock(
		context.Background(),
		"product-1",
		10,
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

func TestProductService_Delete(t *testing.T) {
	var deletedID string

	repository := &mockProductRepository{
		deleteFunc: func(
			ctx context.Context,
			id string,
		) error {
			deletedID = id
			return nil
		},
	}

	service := NewProductService(repository)

	err := service.Delete(
		context.Background(),
		"product-1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if deletedID != "product-1" {
		t.Fatalf(
			"expected ID %q, got %q",
			"product-1",
			deletedID,
		)
	}
}

func TestProductService_Delete_Error(t *testing.T) {
	repository := &mockProductRepository{
		deleteFunc: func(
			ctx context.Context,
			id string,
		) error {
			return errors.New("database error")
		},
	}

	service := NewProductService(repository)

	err := service.Delete(
		context.Background(),
		"product-1",
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
