package service

import (
	"context"

	"github.com/google/uuid"

	"store-service/internal/model"
	"store-service/internal/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, o *model.Order) error {
	return s.repo.Create(ctx, o)
}

func (s *OrderService) Get(ctx context.Context, id uuid.UUID) (model.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *OrderService) List(ctx context.Context, limit, offset int) ([]model.Order, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *OrderService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *OrderService) AddProductToOrder(ctx context.Context, orderID, productID uuid.UUID, qty int) (model.OrderItem, error) {
	return s.repo.AddProductToOrder(ctx, orderID, productID, qty)
}
