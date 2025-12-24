package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"store-service/internal/service"
)

func NewRouter(log *zap.Logger, services *service.Services) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(LoggerMiddleware(log))

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	registerCategoryRoutes(r, services.Categories)
	registerCustomerRoutes(r, services.Customers)
	registerProductRoutes(r, services.Products)
	registerOrderRoutes(r, services.Orders)
	registerReportRoutes(r, services.Reports)
	registerDocsRoutes(r)

	return r
}
