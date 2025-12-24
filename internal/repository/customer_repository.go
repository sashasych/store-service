package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"store-service/internal/model"
)

type CustomerRepository struct {
	pool *pgxpool.Pool
}

func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{pool: pool}
}

func (r *CustomerRepository) Create(ctx context.Context, c *model.Customer) error {
	now := time.Now().UTC()
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = now
	c.UpdatedAt = now

	query := `INSERT INTO customers (id, name, email, phone, address, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query, c.ID, c.Name, c.Email, c.Phone, c.Address, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *CustomerRepository) Get(ctx context.Context, id uuid.UUID) (model.Customer, error) {
	var c model.Customer
	query := `SELECT id, name, email, phone, address, created_at, updated_at FROM customers WHERE id=$1`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c, ErrNotFound
		}
		return c, err
	}
	return c, nil
}

func (r *CustomerRepository) Update(ctx context.Context, c *model.Customer) error {
	c.UpdatedAt = time.Now().UTC()
	query := `UPDATE customers SET name=$1, email=$2, phone=$3, address=$4, updated_at=$5 WHERE id=$6`
	cmd, err := r.pool.Exec(ctx, query, c.Name, c.Email, c.Phone, c.Address, c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM customers WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CustomerRepository) List(ctx context.Context, limit, offset int) ([]model.Customer, error) {
	query := `SELECT id, name, email, phone, address, created_at, updated_at FROM customers ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Address, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}
