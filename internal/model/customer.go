package model

import (
	"time"

	"github.com/google/uuid"
)

// покупатель
// id уникальный идентификатор покупателя
// name имя покупателя
// email email покупателя
// phone телефон покупателя
// address адрес покупателя
// created_at дата создания покупателя
// updated_at дата обновления покупателя
type Customer struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
