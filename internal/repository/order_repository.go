package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"store-service/internal/model"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Create(ctx context.Context, o *model.Order) error {
	now := time.Now().UTC()
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.Status == "" {
		o.Status = "new"
	}
	o.CreatedAt = now
	o.UpdatedAt = now
	if o.TotalPrice.IsZero() {
		o.TotalPrice = decimal.Zero
	}

	query := `INSERT INTO orders (id, customer_id, total_price, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, o.ID, o.CustomerID, o.TotalPrice, o.Status, o.CreatedAt, o.UpdatedAt)
	return err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	cmd, err := r.pool.Exec(ctx, `UPDATE orders SET status=$1, updated_at=$2 WHERE id=$3`, status, time.Now().UTC(), id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *OrderRepository) Get(ctx context.Context, id uuid.UUID) (model.Order, error) {
	var o model.Order

	err := r.pool.QueryRow(ctx, `SELECT id, customer_id, total_price, status, created_at, updated_at FROM orders WHERE id=$1`, id).
		Scan(&o.ID, &o.CustomerID, &o.TotalPrice, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return o, ErrNotFound
		}
		return o, err
	}

	items, err := r.fetchItems(ctx, id)
	if err != nil {
		return o, err
	}
	o.Items = items
	return o, nil
}

func (r *OrderRepository) List(ctx context.Context, limit, offset int) ([]model.Order, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, customer_id, total_price, status, created_at, updated_at FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	var ids []uuid.UUID

	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.TotalPrice, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
		ids = append(ids, o.ID)
	}

	if len(ids) == 0 {
		return orders, rows.Err()
	}

	itemMap, err := r.fetchItemsForOrders(ctx, ids)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		orders[i].Items = itemMap[orders[i].ID]
	}
	return orders, rows.Err()
}

func (r *OrderRepository) fetchItems(ctx context.Context, orderID uuid.UUID) ([]model.OrderItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, order_id, product_id, quantity, sub_total, created_at, updated_at FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItem
	for rows.Next() {
		var it model.OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Quantity, &it.SubTotal, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

func (r *OrderRepository) fetchItemsForOrders(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]model.OrderItem, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, order_id, product_id, quantity, sub_total, created_at, updated_at FROM order_items WHERE order_id = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]model.OrderItem)
	for rows.Next() {
		var it model.OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Quantity, &it.SubTotal, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		result[it.OrderID] = append(result[it.OrderID], it)
	}
	return result, rows.Err()
}

// AddProductToOrder adds or increments a product inside the order with transactional guarantees.
func (r *OrderRepository) AddProductToOrder(ctx context.Context, orderID, productID uuid.UUID, qty int) (model.OrderItem, error) {
	var item model.OrderItem
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return item, err
	}
	defer tx.Rollback(ctx)

	// Ensure order exists and lock it.
	if err := tx.QueryRow(ctx, `SELECT id FROM orders WHERE id=$1 FOR UPDATE`, orderID).Scan(new(uuid.UUID)); err != nil {
		if err == pgx.ErrNoRows {
			return item, ErrNotFound
		}
		return item, err
	}

	var price decimal.Decimal
	var stock int
	if err := tx.QueryRow(ctx, `SELECT price, quantity FROM products WHERE id=$1 FOR UPDATE`, productID).Scan(&price, &stock); err != nil {
		if err == pgx.ErrNoRows {
			return item, ErrNotFound
		}
		return item, err
	}
	if stock < qty {
		return item, ErrNotEnoughStock
	}

	var currentID uuid.UUID
	var currentQty int
	var currentSub decimal.Decimal
	var currentCreated time.Time
	exists := true
	err = tx.QueryRow(ctx, `SELECT id, quantity, sub_total, created_at FROM order_items WHERE order_id=$1 AND product_id=$2 FOR UPDATE`, orderID, productID).
		Scan(&currentID, &currentQty, &currentSub, &currentCreated)
	if err != nil {
		if err == pgx.ErrNoRows {
			exists = false
		} else {
			return item, err
		}
	}

	newQty := qty
	if exists {
		newQty = currentQty + qty
	}
	newSub := price.Mul(decimal.NewFromInt(int64(newQty)))

	now := time.Now().UTC()
	item.OrderID = orderID
	item.ProductID = productID
	item.Quantity = newQty
	item.SubTotal = newSub
	item.UpdatedAt = now

	if exists {
		item.ID = currentID
		item.CreatedAt = currentCreated
		_, err = tx.Exec(ctx, `UPDATE order_items SET quantity=$1, sub_total=$2, updated_at=$3 WHERE order_id=$4 AND product_id=$5`, newQty, newSub, now, orderID, productID)
	} else {
		item.ID = uuid.New()
		_, err = tx.Exec(ctx, `INSERT INTO order_items (id, order_id, product_id, quantity, sub_total, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`, item.ID, orderID, productID, newQty, newSub, now, now)
		item.CreatedAt = now
	}
	if err != nil {
		return item, err
	}

	delta := newSub
	if exists {
		delta = newSub.Sub(currentSub)
	}

	if _, err := tx.Exec(ctx, `UPDATE orders SET total_price = total_price + $1, updated_at=$2 WHERE id=$3`, delta, now, orderID); err != nil {
		return item, err
	}

	if _, err := tx.Exec(ctx, `UPDATE products SET quantity = quantity - $1, updated_at=$2 WHERE id=$3`, qty, now, productID); err != nil {
		return item, err
	}

	if err := tx.Commit(ctx); err != nil {
		return item, err
	}

	item.OrderID = orderID
	item.ProductID = productID
	item.Quantity = newQty
	item.SubTotal = newSub
	item.CreatedAt = now
	item.UpdatedAt = now
	return item, nil
}
