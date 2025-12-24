package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Order struct {
	ID         uuid.UUID       `json:"id"`
	CustomerID uuid.UUID       `json:"customer_id"`
	Items      []OrderItem     `json:"items"`
	TotalPrice decimal.Decimal `json:"total_price"`
	Status     string          `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID       `json:"id"`
	OrderID   uuid.UUID       `json:"order_id"`
	ProductID uuid.UUID       `json:"product_id"`
	Quantity  int             `json:"quantity"`
	SubTotal  decimal.Decimal `json:"sub_total"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
