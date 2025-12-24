package model

import (
	"time"

	"github.com/google/uuid"
)

// связующая таблица для категорий продуктов
// product_id идентификатор продукта
// catagory_id идентификатор категории
// created_at дата создания связи
// updated_at дата обновления связи

// primary key (product_id, catagory_id)

type ProductCatagory struct {
	ProductID  uuid.UUID `json:"product_id"`
	CatagoryID uuid.UUID `json:"catagory_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
