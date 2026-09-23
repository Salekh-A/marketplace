package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"analytics/internal/model"
	"analytics/internal/service"
)

type AnalyticsService interface {
	GetRevenue(rctx context.Context) (model.Revenue, error)
	GetOrders(rctx context.Context) (model.OrderStats, error)
	GetAverageCheck(rctx context.Context) (model.AverageCheck, error)
	GetTopProducts(rctx context.Context) ([]model.ProductStats, error)
}

type AnalyticsHandler struct {
	service AnalyticsService
}

func NewAnalyticsHandler(
	service AnalyticsService,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

func (h *AnalyticsHandler) Revenue(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := h.service.GetRevenue(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get revenue",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *AnalyticsHandler) Orders(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := h.service.GetOrders(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get orders",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *AnalyticsHandler) AverageCheck(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := h.service.GetAverageCheck(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get average check",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *AnalyticsHandler) TopProducts(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := h.service.GetTopProducts(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to get top products",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, result)
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

var _ AnalyticsService = (*service.AnalyticsService)(nil)
