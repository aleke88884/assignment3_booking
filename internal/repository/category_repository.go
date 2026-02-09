package repository

import (
	"context"
	"database/sql"
	"errors"

	"smartbooking/internal/models"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.ResourceCategory) error
	GetByID(ctx context.Context, id int64) (*models.ResourceCategory, error)
	List(ctx context.Context) ([]*models.ResourceCategory, error)
	Update(ctx context.Context, category *models.ResourceCategory) error
	Delete(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.ResourceCategory) error {
	query := `
		INSERT INTO resource_categories (name, slug, description, icon)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.Icon,
	).Scan(&category.ID)

	return err
}

func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*models.ResourceCategory, error) {
	query := `
		SELECT id, name, slug, description, icon
		FROM resource_categories
		WHERE id = $1
	`

	category := &models.ResourceCategory{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.Icon,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCategoryNotFound
	}
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) List(ctx context.Context) ([]*models.ResourceCategory, error) {
	query := `
		SELECT id, name, slug, description, icon
		FROM resource_categories
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*models.ResourceCategory, 0)
	for rows.Next() {
		category := &models.ResourceCategory{}

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Icon,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, rows.Err()
}

func (r *categoryRepository) Update(ctx context.Context, category *models.ResourceCategory) error {
	query := `
		UPDATE resource_categories
		SET name = $1, slug = $2, description = $3, icon = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.Icon,
		category.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM resource_categories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrCategoryNotFound
	}

	return nil
}
