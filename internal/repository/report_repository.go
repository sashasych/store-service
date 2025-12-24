package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type CustomerTotal struct {
	CustomerName string          `json:"customer_name"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
}

type CategoryChildrenCount struct {
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	Count      int       `json:"children_count"`
}

type TopProduct struct {
	ProductName    string `json:"product_name"`
	CategoryLevel1 string `json:"category_level_1"`
	TotalQuantity  int    `json:"total_quantity"`
}

type ReportRepository struct {
	pool *pgxpool.Pool
}

func NewReportRepository(pool *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{pool: pool}
}

func (r *ReportRepository) CustomerTotals(ctx context.Context) ([]CustomerTotal, error) {
	const q = `
SELECT c.name AS customer_name, COALESCE(SUM(oi.sub_total), 0) AS total_amount
FROM customers c
LEFT JOIN orders o ON o.customer_id = c.id
LEFT JOIN order_items oi ON oi.order_id = o.id
GROUP BY c.name
ORDER BY c.name;
`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []CustomerTotal
	for rows.Next() {
		var row CustomerTotal
		if err := rows.Scan(&row.CustomerName, &row.TotalAmount); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, rows.Err()
}

func (r *ReportRepository) CategoryChildren(ctx context.Context) ([]CategoryChildrenCount, error) {
	const q = `
SELECT parent.id, parent.name, COUNT(child.id) AS children_count
FROM categories AS parent
LEFT JOIN categories AS child ON child.parent_id = parent.id
GROUP BY parent.id, parent.name
ORDER BY parent.name;
`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []CategoryChildrenCount
	for rows.Next() {
		var row CategoryChildrenCount
		if err := rows.Scan(&row.CategoryID, &row.Name, &row.Count); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, rows.Err()
}

func (r *ReportRepository) TopProductsLastMonth(ctx context.Context) ([]TopProduct, error) {
	const q = `
WITH RECURSIVE cat_path AS (
    SELECT id, parent_id, name AS root_name, id AS root_id
    FROM categories
    WHERE parent_id IS NULL
    UNION ALL
    SELECT c.id, c.parent_id, cp.root_name, cp.root_id
    FROM categories c
    JOIN cat_path cp ON c.parent_id = cp.id
)
SELECT
    p.name AS product_name,
    cp.root_name AS category_level_1,
    SUM(oi.quantity) AS total_quantity
FROM order_items oi
JOIN orders o ON o.id = oi.order_id
JOIN products p ON p.id = oi.product_id
LEFT JOIN product_catagories pc ON pc.product_id = p.id
LEFT JOIN cat_path cp ON cp.id = pc.catagory_id
WHERE o.created_at >= date_trunc('month', now()) - INTERVAL '1 month'
  AND o.created_at <  date_trunc('month', now())
GROUP BY p.name, cp.root_name
ORDER BY total_quantity DESC
LIMIT 5;
`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []TopProduct
	for rows.Next() {
		var row TopProduct
		if err := rows.Scan(&row.ProductName, &row.CategoryLevel1, &row.TotalQuantity); err != nil {
			return nil, err
		}
		res = append(res, row)
	}
	return res, rows.Err()
}
