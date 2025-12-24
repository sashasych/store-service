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

type productHandler struct {
	svc *service.ProductService
}

func registerProductRoutes(r chi.Router, svc *service.ProductService) {
	h := &productHandler{svc: svc}
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.get)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h *productHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p := req.ToModel(uuid.Nil)

	if err := h.svc.Create(ctx, &p); err != nil {
		log.Error("failed to create product", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to create product")
		return
	}
	writeJSON(w, http.StatusCreated, dto.FromProduct(p))
}

func (h *productHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	p, err := h.svc.Get(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "product not found")
			return
		}
		log.Error("failed to get product", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to get product")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromProduct(p))
}

func (h *productHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var req dto.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p := req.ToModel(id)

	if err := h.svc.Update(ctx, &p); err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "product not found")
			return
		}
		log.Error("failed to update product", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to update product")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromProduct(p))
}

func (h *productHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "product not found")
			return
		}
		log.Error("failed to delete product", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to delete product")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *productHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	limit, offset := parsePagination(r)

	products, err := h.svc.List(ctx, limit, offset)
	if err != nil {
		log.Error("failed to list products", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to list products")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromProducts(products))
}
