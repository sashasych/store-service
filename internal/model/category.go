package model

import (
	"time"

	"github.com/google/uuid"
)

// категория товара
// id уникальный идентификатор категории
// name название категории
// slug уникальный идентификатор категории для url
// parent_id идентификатор родительской категории
// level уровень вложенности категории
// is_active флаг активности категории
// sort_order порядок сортировки категории
// created_at дата создания категории
// updated_at дата обновления категории

type Category struct {
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
