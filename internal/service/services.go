package service

import "store-service/internal/repository"

// Services aggregates all domain services for easier wiring.
type Services struct {
	Categories *CategoryService
	Customers  *CustomerService
	Products   *ProductService
	Orders     *OrderService
	Reports    *ReportService
}

func NewServices(
	categoryRepo *repository.CategoryRepository,
	customerRepo *repository.CustomerRepository,
	productRepo *repository.ProductRepository,
	orderRepo *repository.OrderRepository,
	reportRepo *repository.ReportRepository,
) *Services {
	return &Services{
		Categories: NewCategoryService(categoryRepo),
		Customers:  NewCustomerService(customerRepo),
		Products:   NewProductService(productRepo),
		Orders:     NewOrderService(orderRepo),
		Reports:    NewReportService(reportRepo),
	}
}
