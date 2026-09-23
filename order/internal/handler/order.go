package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"order/internal/model"
	"order/internal/service"
)

type OrderService interface {
	Create(
		ctx context.Context,
		productID string,
		quantity int,
		price float64,
	) (model.Order, error)

	GetByID(
		ctx context.Context,
		id string,
	) (model.Order, error)

	Cancel(
		ctx context.Context,
		id string,
	) error

	Delete(
		ctx context.Context,
		id string,
	) error
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

type createOrderRequest struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

func (h *OrderHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request createOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	order, err := h.service.Create(
		r.Context(),
		request.ProductID,
		request.Quantity,
		request.Price,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")

	if id == "" {
		http.Error(
			w,
			"order id is required",
			http.StatusBadRequest,
		)
		return
	}

	order, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			"order not found",
			http.StatusNotFound,
		)
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) Cancel(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")

	if id == "" {
		http.Error(
			w,
			"order id is required",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Cancel(
		r.Context(),
		id,
	); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")

	if id == "" {
		http.Error(
			w,
			"order id is required",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Delete(
		r.Context(),
		id,
	); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

var _ OrderService = (*service.OrderService)(nil)
