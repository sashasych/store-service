package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	Quantity  int             `json:"quantity"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
