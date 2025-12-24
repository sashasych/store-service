package application

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"store-service/internal/api"
	"store-service/internal/config"
	"store-service/internal/logger"
	"store-service/internal/repository"
	"store-service/internal/service"
)

type Application struct {
	cfg    config.Config
	log    *zap.Logger
	db     *pgxpool.Pool
	router *chi.Mux
	server *http.Server
}

// New builds application with all dependencies.
func New(ctx context.Context) (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	dbCfg, err := pgxpool.ParseConfig(cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}
	dbCfg.MaxConns = cfg.Postgres.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		return nil, err
	}

	categoryRepo := repository.NewCategoryRepository(pool)
	customerRepo := repository.NewCustomerRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)
	reportRepo := repository.NewReportRepository(pool)

	services := service.NewServices(categoryRepo, customerRepo, productRepo, orderRepo, reportRepo)
	router := api.NewRouter(log, services)

	server := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Application{
		cfg:    cfg,
		log:    log,
		db:     pool,
		router: router,
		server: server,
	}, nil
}

// Run starts HTTP server and waits for shutdown signal.
func (a *Application) Run(ctx context.Context) error {
	ctx = logger.WithContext(ctx, a.log)
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.log.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			return err
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.GracefulTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
		return err
	}

	a.db.Close()
	return nil
}
