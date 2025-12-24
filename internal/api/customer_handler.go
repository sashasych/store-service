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

type customerHandler struct {
	svc *service.CustomerService
}

func registerCustomerRoutes(r chi.Router, svc *service.CustomerService) {
	h := &customerHandler{svc: svc}
	r.Route("/customers", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.get)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h *customerHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req dto.CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c := req.ToModel(uuid.Nil)

	if err := h.svc.Create(ctx, &c); err != nil {
		log.Error("failed to create customer", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to create customer")
		return
	}
	writeJSON(w, http.StatusCreated, dto.FromCustomer(c))
}

func (h *customerHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	c, err := h.svc.Get(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "customer not found")
			return
		}
		log.Error("failed to get customer", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to get customer")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromCustomer(c))
}

func (h *customerHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	var req dto.CustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c := req.ToModel(id)

	if err := h.svc.Update(ctx, &c); err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "customer not found")
			return
		}
		log.Error("failed to update customer", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to update customer")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromCustomer(c))
}

func (h *customerHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "customer not found")
			return
		}
		log.Error("failed to delete customer", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to delete customer")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *customerHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	limit, offset := parsePagination(r)

	customers, err := h.svc.List(ctx, limit, offset)
	if err != nil {
		log.Error("failed to list customers", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to list customers")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromCustomers(customers))
}
