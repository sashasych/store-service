package service

import (
	"context"

	"github.com/google/uuid"

	"store-service/internal/model"
	"store-service/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) error {
	return s.repo.Create(ctx, p)
}

func (s *ProductService) Get(ctx context.Context, id uuid.UUID) (model.Product, error) {
	return s.repo.Get(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, p *model.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) List(ctx context.Context, limit, offset int) ([]model.Product, error) {
	return s.repo.List(ctx, limit, offset)
}
