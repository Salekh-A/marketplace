package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"product/internal/model"
	"product/internal/service"
)

type ProductService interface {
	Create(
		ctx context.Context,
		name string,
		price float64,
		stock int,
	) (model.Product, error)

	GetAll(
		ctx context.Context,
	) ([]model.Product, error)

	GetByID(
		ctx context.Context,
		id string,
	) (model.Product, error)

	Update(
		ctx context.Context,
		id string,
		name string,
		price float64,
		stock int,
	) (model.Product, error)

	Delete(
		ctx context.Context,
		id string,
	) error
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

type productRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func (h *ProductHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request productRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	product, err := h.service.Create(
		r.Context(),
		request.Name,
		request.Price,
		request.Stock,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	writeJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	products, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get products",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	if id == "" {
		http.Error(
			w,
			"product id is required",
			http.StatusBadRequest,
		)
		return
	}

	product, err := h.service.GetByID(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			"product not found",
			http.StatusNotFound,
		)
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	if id == "" {
		http.Error(
			w,
			"product id is required",
			http.StatusBadRequest,
		)
		return
	}

	var request productRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	product, err := h.service.Update(
		r.Context(),
		id,
		request.Name,
		request.Price,
		request.Stock,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	if id == "" {
		http.Error(
			w,
			"product id is required",
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}

// compile-time check
var _ ProductService = (*service.ProductService)(nil)
