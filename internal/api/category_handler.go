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

type categoryHandler struct {
	svc *service.CategoryService
}

func registerCategoryRoutes(r chi.Router, svc *service.CategoryService) {
	h := &categoryHandler{svc: svc}
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.get)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h *categoryHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var req dto.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c := req.ToModel(uuid.Nil)

	if err := h.svc.Create(ctx, &c); err != nil {
		log.Error("failed to create category", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to create category")
		return
	}
	writeJSON(w, http.StatusCreated, dto.FromCategory(c))
}

func (h *categoryHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	c, err := h.svc.Get(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "category not found")
			return
		}
		log.Error("failed to get category", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to get category")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromCategory(c))
}

func (h *categoryHandler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	var req dto.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c := req.ToModel(id)

	if err := h.svc.Update(ctx, &c); err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "category not found")
			return
		}
		log.Error("failed to update category", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to update category")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromCategory(c))
}

func (h *categoryHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if err == repository.ErrNotFound {
			writeError(w, http.StatusNotFound, "category not found")
			return
		}
		log.Error("failed to delete category", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to delete category")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *categoryHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	limit, offset := parsePagination(r)

	categories, err := h.svc.List(ctx, limit, offset)
	if err != nil {
		log.Error("failed to list categories", zapError(err))
		writeError(w, http.StatusInternalServerError, "failed to list categories")
		return
	}
	writeJSON(w, http.StatusOK, dto.FromCategories(categories))
}
