package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"store-service/internal/model"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) error {
	now := time.Now().UTC()
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	p.CreatedAt = now
	p.UpdatedAt = now
	query := `INSERT INTO products (id, name, price, quantity, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, p.ID, p.Name, p.Price, p.Quantity, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *ProductRepository) Get(ctx context.Context, id uuid.UUID) (model.Product, error) {
	var p model.Product
	query := `SELECT id, name, price, quantity, created_at, updated_at FROM products WHERE id=$1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Quantity, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return p, ErrNotFound
		}
		return p, err
	}
	return p, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *model.Product) error {
	p.UpdatedAt = time.Now().UTC()
	query := `UPDATE products SET name=$1, price=$2, quantity=$3, updated_at=$4 WHERE id=$5`
	cmd, err := r.pool.Exec(ctx, query, p.Name, p.Price, p.Quantity, p.UpdatedAt, p.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]model.Product, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, price, quantity, created_at, updated_at FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}
