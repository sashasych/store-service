package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"store-service/internal/model"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.Category) error {
	now := time.Now().UTC()
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.CreatedAt = now
	c.UpdatedAt = now

	query := `INSERT INTO categories 
		(id, name, slug, parent_id, level, is_active, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.pool.Exec(ctx, query, c.ID, c.Name, c.Slug, c.ParentID, c.Level, c.IsActive, c.SortOrder, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *CategoryRepository) Get(ctx context.Context, id uuid.UUID) (model.Category, error) {
	var c model.Category

	query := `SELECT id, name, slug, parent_id, level, is_active, sort_order, created_at, updated_at 
		FROM categories WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.ParentID, &c.Level, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c, ErrNotFound
		}
		return c, err
	}
	return c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *model.Category) error {
	c.UpdatedAt = time.Now().UTC()
	query := `UPDATE categories SET name=$1, slug=$2, parent_id=$3, level=$4, is_active=$5, sort_order=$6, updated_at=$7 WHERE id=$8`
	cmd, err := r.pool.Exec(ctx, query, c.Name, c.Slug, c.ParentID, c.Level, c.IsActive, c.SortOrder, c.UpdatedAt, c.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM categories WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *CategoryRepository) List(ctx context.Context, limit, offset int) ([]model.Category, error) {
	query := `SELECT id, name, slug, parent_id, level, is_active, sort_order, created_at, updated_at 
		FROM categories ORDER BY sort_order ASC, created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.ParentID, &c.Level, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}
