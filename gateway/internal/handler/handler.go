package handler

import (
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"

	"gateway/internal/grpc"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	orderClient     *grpc.OrderClient
	productClient   *grpc.ProductClient
	analyticsClient *grpc.AnalyticsClient
}

func NewHandler(
	orderClient *grpc.OrderClient,
	productClient *grpc.ProductClient,
	analyticsClient *grpc.AnalyticsClient,
) *Handler {
	return &Handler{
		orderClient:     orderClient,
		productClient:   productClient,
		analyticsClient: analyticsClient,
	}
}

func (h *Handler) GetOrder(
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

	order, err := h.orderClient.GetOrder(
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

func (h *Handler) GetProduct(
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

	product, err := h.productClient.GetProduct(
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

func (h *Handler) GetProducts(
	w http.ResponseWriter,
	r *http.Request,
) {
	products, err := h.productClient.GetProducts(
		r.Context(),
	)
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

func (h *Handler) GetRevenue(
	w http.ResponseWriter,
	r *http.Request,
) {
	revenue, err := h.analyticsClient.GetRevenue(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to get revenue",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, revenue)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ProductID string  `json:"product_id"`
		Quantity  int32   `json:"quantity"`
		Price     float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.orderClient.CreateOrder(
		r.Context(),
		request.ProductID,
		request.Quantity,
		request.Price,
	)
	if err != nil {
		code := http.StatusInternalServerError

		if status.Code(err) == codes.InvalidArgument {
			code = http.StatusBadRequest
		}

		http.Error(w, err.Error(), code)
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func (h *Handler) GetOrders(
	w http.ResponseWriter,
	r *http.Request,
) {
	orders, err := h.analyticsClient.GetOrders(
		r.Context(),
	)
	if err != nil {
		http.Error(
			w,
			"failed to get orders",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

func (h *Handler) CancelOrder(
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

	response, err := h.orderClient.CancelOrder(
		r.Context(),
		id,
	)
	if err != nil {
		http.Error(
			w,
			"failed to cancel order",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
) {
	w.Header().Set("Content-Type", "application/json")

	if message, ok := data.(proto.Message); ok {
		body, err := protojson.MarshalOptions{
			EmitUnpopulated: true,
		}.Marshal(message)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(status)
		w.Write(body)
		w.Write([]byte("\n"))

		return
	}

	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
	}
}
