package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"product/internal/model"
)

func TestProductHandler_Create_Success(t *testing.T) {
	expected := model.Product{
		ID:        "product-1",
		Name:      "iPhone",
		Price:     999.99,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockProductService{
		createFunc: func(
			ctx context.Context,
			name string,
			price float64,
			stock int,
		) (model.Product, error) {
			if name != "iPhone" {
				t.Errorf("expected name iPhone, got %s", name)
			}

			if price != 999.99 {
				t.Errorf("expected price 999.99, got %v", price)
			}

			if stock != 10 {
				t.Errorf("expected stock 10, got %d", stock)
			}

			return expected, nil
		},
	}

	handler := NewProductHandler(mockService)

	body := `{
		"name": "iPhone",
		"price": 999.99,
		"stock": 10
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			rec.Code,
		)
	}

	var response model.Product

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != expected.ID {
		t.Errorf("expected id %q, got %q", expected.ID, response.ID)
	}

	if response.Name != expected.Name {
		t.Errorf("expected name %q, got %q", expected.Name, response.Name)
	}
}

func TestProductHandler_Create_InvalidJSON(t *testing.T) {
	handler := NewProductHandler(&mockProductService{})

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		strings.NewReader(`invalid json`),
	)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestProductHandler_Create_ServiceError(t *testing.T) {
	mockService := &mockProductService{
		createFunc: func(
			ctx context.Context,
			name string,
			price float64,
			stock int,
		) (model.Product, error) {
			return model.Product{}, errors.New("name is required")
		},
	}

	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		strings.NewReader(`{
			"name": "",
			"price": 100,
			"stock": 10
		}`),
	)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestProductHandler_GetAll_Success(t *testing.T) {
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

	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response []model.Product

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != len(expected) {
		t.Fatalf(
			"expected %d products, got %d",
			len(expected),
			len(response),
		)
	}

	if response[0].ID != expected[0].ID {
		t.Errorf(
			"expected first product id %q, got %q",
			expected[0].ID,
			response[0].ID,
		)
	}
}

func TestProductHandler_GetAll_ServiceError(t *testing.T) {
	mockService := &mockProductService{
		getAllFunc: func(
			ctx context.Context,
		) ([]model.Product, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}

func TestProductHandler_GetByID_Success(t *testing.T) {
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
				t.Errorf("expected id product-1, got %s", id)
			}

			return expected, nil
		},
	}

	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/product-1",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response model.Product

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != expected.ID {
		t.Errorf("expected id %q, got %q", expected.ID, response.ID)
	}
}

func TestProductHandler_GetByID_EmptyID(t *testing.T) {
	handler := NewProductHandler(&mockProductService{})

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestProductHandler_GetByID_NotFound(t *testing.T) {
	mockService := &mockProductService{
		getByIDFunc: func(
			ctx context.Context,
			id string,
		) (model.Product, error) {
			return model.Product{}, errors.New("product not found")
		},
	}

	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/product-1",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			rec.Code,
		)
	}
}

func TestProductHandler_Update_Success(t *testing.T) {
	expected := model.Product{
		ID:    "product-1",
		Name:  "Updated iPhone",
		Price: 1099.99,
		Stock: 20,
	}

	mockService := &mockProductService{
		updateFunc: func(
			ctx context.Context,
			id string,
			name string,
			price float64,
			stock int,
		) (model.Product, error) {
			if id != "product-1" {
				t.Errorf("expected id product-1, got %s", id)
			}

			if name != "Updated iPhone" {
				t.Errorf(
					"expected name Updated iPhone, got %s",
					name,
				)
			}

			if price != 1099.99 {
				t.Errorf(
					"expected price 1099.99, got %v",
					price,
				)
			}

			if stock != 20 {
				t.Errorf(
					"expected stock 20, got %d",
					stock,
				)
			}

			return expected, nil
		},
	}

	handler := NewProductHandler(mockService)

	body := `{
		"name": "Updated iPhone",
		"price": 1099.99,
		"stock": 20
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/product-1",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response model.Product

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Name != expected.Name {
		t.Errorf(
			"expected name %q, got %q",
			expected.Name,
			response.Name,
		)
	}
}

func TestProductHandler_Update_EmptyID(t *testing.T) {
	handler := NewProductHandler(&mockProductService{})

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/",
		strings.NewReader(`{
			"name": "iPhone",
			"price": 100,
			"stock": 10
		}`),
	)
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestProductHandler_Update_InvalidJSON(t *testing.T) {
	handler := NewProductHandler(&mockProductService{})

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/product-1",
		strings.NewReader(`invalid json`),
	)
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestProductHandler_Update_ServiceError(t *testing.T) {
	mockService := &mockProductService{
		updateFunc: func(
			ctx context.Context,
			id string,
			name string,
			price float64,
			stock int,
		) (model.Product, error) {
			return model.Product{}, errors.New("price cannot be negative")
		},
	}

	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPut,
		"/products/product-1",
		strings.NewReader(`{
			"name": "iPhone",
			"price": -100,
			"stock": 10
		}`),
	)
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestProductHandler_Delete_Success(t *testing.T) {
	mockService := &mockProductService{
		deleteFunc: func(
			ctx context.Context,
			id string,
		) error {
			if id != "product-1" {
				t.Errorf("expected id product-1, got %s", id)
			}

			return nil
		},
	}

	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/products/product-1",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			rec.Code,
		)
	}
}

func TestProductHandler_Delete_EmptyID(t *testing.T) {
	handler := NewProductHandler(&mockProductService{})

	req := httptest.NewRequest(
		http.MethodDelete,
		"/products/",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestProductHandler_Delete_ServiceError(t *testing.T) {
	mockService := &mockProductService{
		deleteFunc: func(
			ctx context.Context,
			id string,
		) error {
			return errors.New("database error")
		},
	}

	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/products/product-1",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.Delete(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}
