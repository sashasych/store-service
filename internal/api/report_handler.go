package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"store-service/internal/api/dto"
	"store-service/internal/logger"
	"store-service/internal/service"
)

type reportHandler struct {
	svc *service.ReportService
}

func registerReportRoutes(r chi.Router, svc *service.ReportService) {
	h := &reportHandler{svc: svc}
	r.Route("/reports", func(r chi.Router) {
		r.Get("/customer-totals", h.customerTotals)
		r.Get("/category-children", h.categoryChildren)
		r.Get("/top-products-last-month", h.topProductsLastMonth)
	})
}

func (h *reportHandler) customerTotals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	data, err := h.svc.CustomerTotals(ctx)
	if err != nil {
		log.Error("failed to fetch customer totals", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to fetch customer totals")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromCustomerTotals(data))
}

func (h *reportHandler) categoryChildren(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	data, err := h.svc.CategoryChildren(ctx)
	if err != nil {
		log.Error("failed to fetch category children", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to fetch category children")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromCategoryChildren(data))
}

func (h *reportHandler) topProductsLastMonth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	data, err := h.svc.TopProductsLastMonth(ctx)
	if err != nil {
		log.Error("failed to fetch top products", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to fetch top products")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromTopProducts(data))
}
