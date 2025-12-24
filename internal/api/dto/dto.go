package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"store-service/internal/model"
	"store-service/internal/repository"
)

// Category DTOs
type CategoryRequest struct {
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Level     int        `json:"level"`
	IsActive  bool       `json:"is_active"`
	SortOrder int        `json:"sort_order"`
}

type CategoryResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	Level     int        `json:"level"`
	IsActive  bool       `json:"is_active"`
	SortOrder int        `json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (r CategoryRequest) ToModel(id uuid.UUID) model.Category {
	return model.Category{
		ID:        id,
		Name:      r.Name,
		Slug:      r.Slug,
		ParentID:  r.ParentID,
		Level:     r.Level,
		IsActive:  r.IsActive,
		SortOrder: r.SortOrder,
	}
}

func FromCategory(m model.Category) CategoryResponse {
	return CategoryResponse{
		ID:        m.ID,
		Name:      m.Name,
		Slug:      m.Slug,
		ParentID:  m.ParentID,
		Level:     m.Level,
		IsActive:  m.IsActive,
		SortOrder: m.SortOrder,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromCategories(list []model.Category) []CategoryResponse {
	result := make([]CategoryResponse, 0, len(list))
	for _, c := range list {
		result = append(result, FromCategory(c))
	}
	return result
}

// Customer DTOs
type CustomerRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type CustomerResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r CustomerRequest) ToModel(id uuid.UUID) model.Customer {
	return model.Customer{
		ID:      id,
		Name:    r.Name,
		Email:   r.Email,
		Phone:   r.Phone,
		Address: r.Address,
	}
}

func FromCustomer(m model.Customer) CustomerResponse {
	return CustomerResponse{
		ID:        m.ID,
		Name:      m.Name,
		Email:     m.Email,
		Phone:     m.Phone,
		Address:   m.Address,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromCustomers(list []model.Customer) []CustomerResponse {
	result := make([]CustomerResponse, 0, len(list))
	for _, c := range list {
		result = append(result, FromCustomer(c))
	}
	return result
}

// Product DTOs
type ProductRequest struct {
	Name     string          `json:"name"`
	Price    decimal.Decimal `json:"price"`
	Quantity int             `json:"quantity"`
}

type ProductResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	Quantity  int             `json:"quantity"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (r ProductRequest) ToModel(id uuid.UUID) model.Product {
	return model.Product{
		ID:       id,
		Name:     r.Name,
		Price:    r.Price,
		Quantity: r.Quantity,
	}
}

func FromProduct(m model.Product) ProductResponse {
	return ProductResponse{
		ID:        m.ID,
		Name:      m.Name,
		Price:     m.Price,
		Quantity:  m.Quantity,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromProducts(list []model.Product) []ProductResponse {
	result := make([]ProductResponse, 0, len(list))
	for _, p := range list {
		result = append(result, FromProduct(p))
	}
	return result
}

// Order DTOs
type OrderRequest struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Status     string    `json:"status"`
}

func (r OrderRequest) ToModel(id uuid.UUID) model.Order {
	return model.Order{
		ID:         id,
		CustomerID: r.CustomerID,
		Status:     r.Status,
	}
}

type OrderItemResponse struct {
	ID        uuid.UUID       `json:"id"`
	ProductID uuid.UUID       `json:"product_id"`
	Quantity  int             `json:"quantity"`
	SubTotal  decimal.Decimal `json:"sub_total"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type OrderResponse struct {
	ID         uuid.UUID           `json:"id"`
	CustomerID uuid.UUID           `json:"customer_id"`
	Items      []OrderItemResponse `json:"items"`
	TotalPrice decimal.Decimal     `json:"total_price"`
	Status     string              `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

func FromOrder(m model.Order) OrderResponse {
	items := make([]OrderItemResponse, 0, len(m.Items))
	for _, it := range m.Items {
		items = append(items, OrderItemResponse{
			ID:        it.ID,
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			SubTotal:  it.SubTotal,
			CreatedAt: it.CreatedAt,
			UpdatedAt: it.UpdatedAt,
		})
	}

	return OrderResponse{
		ID:         m.ID,
		CustomerID: m.CustomerID,
		Items:      items,
		TotalPrice: m.TotalPrice,
		Status:     m.Status,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func FromOrders(list []model.Order) []OrderResponse {
	result := make([]OrderResponse, 0, len(list))
	for _, o := range list {
		result = append(result, FromOrder(o))
	}
	return result
}

type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

// Report DTOs
type CustomerTotalResponse struct {
	CustomerName string          `json:"customer_name"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
}

type CategoryChildrenResponse struct {
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	Count      int       `json:"children_count"`
}

type TopProductResponse struct {
	ProductName    string `json:"product_name"`
	CategoryLevel1 string `json:"category_level_1"`
	TotalQuantity  int    `json:"total_quantity"`
}

func FromCustomerTotals(list []repository.CustomerTotal) []CustomerTotalResponse {
	result := make([]CustomerTotalResponse, 0, len(list))
	for _, row := range list {
		result = append(result, CustomerTotalResponse{
			CustomerName: row.CustomerName,
			TotalAmount:  row.TotalAmount,
		})
	}
	return result
}

func FromCategoryChildren(list []repository.CategoryChildrenCount) []CategoryChildrenResponse {
	result := make([]CategoryChildrenResponse, 0, len(list))
	for _, row := range list {
		result = append(result, CategoryChildrenResponse{
			CategoryID: row.CategoryID,
			Name:       row.Name,
			Count:      row.Count,
		})
	}
	return result
}

func FromTopProducts(list []repository.TopProduct) []TopProductResponse {
	result := make([]TopProductResponse, 0, len(list))
	for _, row := range list {
		result = append(result, TopProductResponse{
			ProductName:    row.ProductName,
			CategoryLevel1: row.CategoryLevel1,
			TotalQuantity:  row.TotalQuantity,
		})
	}
	return result
}
