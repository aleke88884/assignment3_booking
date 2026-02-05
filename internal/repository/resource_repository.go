package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"smartbooking/internal/models"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
)

// ResourceRepository defines the interface for resource data operations
type ResourceRepository interface {
	Create(ctx context.Context, resource *models.Resource) error
	GetByID(ctx context.Context, id int64) (*models.Resource, error)
	Update(ctx context.Context, resource *models.Resource) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.Resource, error)
}

// resourceRepository implements ResourceRepository interface with PostgreSQL storage
type resourceRepository struct {
	db *sql.DB
}

// NewResourceRepository creates a new instance of ResourceRepository
func NewResourceRepository(db *sql.DB) ResourceRepository {
	return &resourceRepository{
		db: db,
	}
}

func (r *resourceRepository) Create(ctx context.Context, resource *models.Resource) error {
	query := `
		INSERT INTO resources (name, description, capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	resource.CreatedAt = now
	resource.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		resource.Name,
		resource.Description,
		resource.Capacity,
		resource.CreatedAt,
		resource.UpdatedAt,
	).Scan(&resource.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *resourceRepository) GetByID(ctx context.Context, id int64) (*models.Resource, error) {
	query := `
		SELECT id, name, description, capacity, created_at, updated_at
		FROM resources
		WHERE id = $1
	`

	resource := &models.Resource{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&resource.ID,
		&resource.Name,
		&resource.Description,
		&resource.Capacity,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrResourceNotFound
	}
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *resourceRepository) Update(ctx context.Context, resource *models.Resource) error {
	query := `
		UPDATE resources
		SET name = $1, description = $2, capacity = $3, updated_at = $4
		WHERE id = $5
	`

	resource.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		resource.Name,
		resource.Description,
		resource.Capacity,
		resource.UpdatedAt,
		resource.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrResourceNotFound
	}

	return nil
}

func (r *resourceRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM resources WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrResourceNotFound
	}

	return nil
}

func (r *resourceRepository) List(ctx context.Context) ([]*models.Resource, error) {
	query := `
		SELECT id, name, description, capacity, created_at, updated_at
		FROM resources
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resources := make([]*models.Resource, 0)
	for rows.Next() {
		resource := &models.Resource{}
		err := rows.Scan(
			&resource.ID,
			&resource.Name,
			&resource.Description,
			&resource.Capacity,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}
