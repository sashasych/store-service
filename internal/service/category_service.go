package service

import (
	"context"

	"github.com/google/uuid"

	"store-service/internal/model"
	"store-service/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, c *model.Category) error {
	return s.repo.Create(ctx, c)
}

func (s *CategoryService) Get(ctx context.Context, id uuid.UUID) (model.Category, error) {
	return s.repo.Get(ctx, id)
}

func (s *CategoryService) Update(ctx context.Context, c *model.Category) error {
	return s.repo.Update(ctx, c)
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *CategoryService) List(ctx context.Context, limit, offset int) ([]model.Category, error) {
	return s.repo.List(ctx, limit, offset)
}
