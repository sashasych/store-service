package service

import (
	"context"

	"github.com/google/uuid"

	"store-service/internal/model"
	"store-service/internal/repository"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(ctx context.Context, c *model.Customer) error {
	return s.repo.Create(ctx, c)
}

func (s *CustomerService) Get(ctx context.Context, id uuid.UUID) (model.Customer, error) {
	return s.repo.Get(ctx, id)
}

func (s *CustomerService) Update(ctx context.Context, c *model.Customer) error {
	return s.repo.Update(ctx, c)
}

func (s *CustomerService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *CustomerService) List(ctx context.Context, limit, offset int) ([]model.Customer, error) {
	return s.repo.List(ctx, limit, offset)
}
