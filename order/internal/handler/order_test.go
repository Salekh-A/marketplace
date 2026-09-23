package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"order/internal/model"
)

func TestOrderHandler_Create(t *testing.T) {
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
			createFunc: func(
				ctx context.Context,
				productID string,
				quantity int,
				price float64,
			) (model.Order, error) {
				if productID != "product-1" {
					t.Errorf(
						"expected product_id product-1, got %q",
						productID,
					)
				}

				if quantity != 2 {
					t.Errorf(
						"expected quantity 2, got %d",
						quantity,
					)
				}

				if price != 999.99 {
					t.Errorf(
						"expected price 999.99, got %v",
						price,
					)
				}

				return order, nil
			},
		}

		handler := NewOrderHandler(service)

		body := `{
			"product_id": "product-1",
			"quantity": 2,
			"price": 999.99
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/orders",
			bytes.NewBufferString(body),
		)

		recorder := httptest.NewRecorder()

		handler.Create(recorder, req)

		if recorder.Code != http.StatusCreated {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusCreated,
				recorder.Code,
			)
		}

		var response model.Order

		if err := json.NewDecoder(
			recorder.Body,
		).Decode(&response); err != nil {
			t.Fatalf(
				"failed to decode response: %v",
				err,
			)
		}

		if response.ID != order.ID {
			t.Errorf(
				"expected id %q, got %q",
				order.ID,
				response.ID,
			)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		handler := NewOrderHandler(
			&mockOrderService{},
		)

		req := httptest.NewRequest(
			http.MethodPost,
			"/orders",
			bytes.NewBufferString(`invalid json`),
		)

		recorder := httptest.NewRecorder()

		handler.Create(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		expectedErr := errors.New("product not found")

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

		handler := NewOrderHandler(service)

		body := `{
			"product_id": "product-1",
			"quantity": 2,
			"price": 999.99
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/orders",
			bytes.NewBufferString(body),
		)

		recorder := httptest.NewRecorder()

		handler.Create(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}
	})
}

func TestOrderHandler_GetByID(t *testing.T) {
	order := model.Order{
		ID:        "order-1",
		ProductID: "product-1",
		Quantity:  2,
		Price:     999.99,
		Status:    "created",
	}

	t.Run("success", func(t *testing.T) {
		service := &mockOrderService{
			getByIDFunc: func(
				ctx context.Context,
				id string,
			) (model.Order, error) {
				if id != order.ID {
					t.Errorf(
						"expected id %q, got %q",
						order.ID,
						id,
					)
				}

				return order, nil
			},
		}

		handler := NewOrderHandler(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/orders/order-1",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.GetByID(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				recorder.Code,
			)
		}

		var response model.Order

		if err := json.NewDecoder(
			recorder.Body,
		).Decode(&response); err != nil {
			t.Fatalf(
				"failed to decode response: %v",
				err,
			)
		}

		if response.ID != order.ID {
			t.Errorf(
				"expected id %q, got %q",
				order.ID,
				response.ID,
			)
		}
	})

	t.Run("empty id", func(t *testing.T) {
		handler := NewOrderHandler(
			&mockOrderService{},
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/orders/",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.GetByID(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}
	})

	t.Run("not found", func(t *testing.T) {
		service := &mockOrderService{
			getByIDFunc: func(
				ctx context.Context,
				id string,
			) (model.Order, error) {
				return model.Order{}, errors.New(
					"order not found",
				)
			},
		}

		handler := NewOrderHandler(service)

		req := httptest.NewRequest(
			http.MethodGet,
			"/orders/order-1",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.GetByID(recorder, req)

		if recorder.Code != http.StatusNotFound {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNotFound,
				recorder.Code,
			)
		}
	})
}

func TestOrderHandler_Cancel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mockOrderService{
			cancelFunc: func(
				ctx context.Context,
				id string,
			) error {
				if id != "order-1" {
					t.Errorf(
						"expected id order-1, got %q",
						id,
					)
				}

				return nil
			},
		}

		handler := NewOrderHandler(service)

		req := httptest.NewRequest(
			http.MethodPost,
			"/orders/order-1/cancel",
			nil,
		)

		req.URL.Path = "/orders/order-1"

		recorder := httptest.NewRecorder()

		handler.Cancel(recorder, req)

		if recorder.Code != http.StatusNoContent {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNoContent,
				recorder.Code,
			)
		}
	})

	t.Run("empty id", func(t *testing.T) {
		handler := NewOrderHandler(
			&mockOrderService{},
		)

		req := httptest.NewRequest(
			http.MethodPost,
			"/orders/",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.Cancel(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service := &mockOrderService{
			cancelFunc: func(
				ctx context.Context,
				id string,
			) error {
				return errors.New("failed to cancel order")
			},
		}

		handler := NewOrderHandler(service)

		req := httptest.NewRequest(
			http.MethodPost,
			"/orders/order-1",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.Cancel(recorder, req)

		if recorder.Code != http.StatusInternalServerError {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusInternalServerError,
				recorder.Code,
			)
		}
	})
}

func TestOrderHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mockOrderService{
			deleteFunc: func(
				ctx context.Context,
				id string,
			) error {
				if id != "order-1" {
					t.Errorf(
						"expected id order-1, got %q",
						id,
					)
				}

				return nil
			},
		}

		handler := NewOrderHandler(service)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/orders/order-1",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.Delete(recorder, req)

		if recorder.Code != http.StatusNoContent {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusNoContent,
				recorder.Code,
			)
		}
	})

	t.Run("empty id", func(t *testing.T) {
		handler := NewOrderHandler(
			&mockOrderService{},
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/orders/",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.Delete(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusBadRequest,
				recorder.Code,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		service := &mockOrderService{
			deleteFunc: func(
				ctx context.Context,
				id string,
			) error {
				return errors.New("failed to delete order")
			},
		}

		handler := NewOrderHandler(service)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/orders/order-1",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.Delete(recorder, req)

		if recorder.Code != http.StatusInternalServerError {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusInternalServerError,
				recorder.Code,
			)
		}
	})
}
