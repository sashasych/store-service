package service

import (
	"context"

	"store-service/internal/repository"
)

type ReportService struct {
	repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) CustomerTotals(ctx context.Context) ([]repository.CustomerTotal, error) {
	return s.repo.CustomerTotals(ctx)
}

func (s *ReportService) CategoryChildren(ctx context.Context) ([]repository.CategoryChildrenCount, error) {
	return s.repo.CategoryChildren(ctx)
}

func (s *ReportService) TopProductsLastMonth(ctx context.Context) ([]repository.TopProduct, error) {
	return s.repo.TopProductsLastMonth(ctx)
}
