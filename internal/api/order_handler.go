package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"store-service/internal/api/dto"
	"store-service/internal/logger"
	"store-service/internal/repository"
	"store-service/internal/service"
)

type orderHandler struct {
	svc *service.OrderService
}

func registerOrderRoutes(r chi.Router, svc *service.OrderService) {
	h := &orderHandler{svc: svc}
	r.Route("/orders", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.get)
		r.Put("/{id}", h.updateStatus)
		r.Delete("/{id}", h.delete)
		r.Post("/{id}/items", h.addItem)
	})
}

func (h *orderHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req dto.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	o := req.ToModel(uuid.Nil)

	if err := h.svc.Create(ctx, &o); err != nil {
		log.Error("failed to create order", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to create order")
		return
	}
	writeJSON(w, http.StatusCreated, dto.FromOrder(o))
}

func (h *orderHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	o, err := h.svc.Get(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "order not found")
			return
		}
		log.Error("failed to get order", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to get order")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromOrder(o))
}

func (h *orderHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	limit, offset := parsePagination(r)

	orders, err := h.svc.List(ctx, limit, offset)
	if err != nil {
		log.Error("failed to list orders", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to list orders")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromOrders(orders))
}

func (h *orderHandler) updateStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	var req dto.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.UpdateStatus(ctx, id, req.Status); err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "order not found")
			return
		}
		log.Error("failed to update order", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to update order")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": req.Status})
}

func (h *orderHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "order not found")
			return
		}
		log.Error("failed to delete order", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to delete order")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *orderHandler) addItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	orderID, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	var req dto.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "quantity must be positive")
		return
	}

	item, err := h.svc.AddProductToOrder(ctx, orderID, req.ProductID, req.Quantity)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			writeError(w, http.StatusNotFound, "order or product not found")
			return
		case repository.ErrNotEnoughStock:
			writeError(w, http.StatusBadRequest, "not enough stock")
			return
		default:
			log.Error("failed to add item to order", zapError(err))
			writeError(w, http.StatusInternalServerError, "failed to add item to order")
			return
		}
	}
	writeJSON(w, http.StatusOK, item)
}
