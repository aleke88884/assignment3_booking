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
		INSERT INTO resources (name, description, capacity, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	resource.CreatedAt = now
	resource.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		resource.Name,
		resource.Description,
		resource.Capacity,
		resource.OwnerID,
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
		SELECT r.id, r.name, r.description, r.capacity, r.owner_id, r.created_at, r.updated_at, u.name as owner_name
		FROM resources r
		LEFT JOIN users u ON r.owner_id = u.id
		WHERE r.id = $1
	`

	resource := &models.Resource{}
	var ownerID sql.NullInt64
	var ownerName sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&resource.ID,
		&resource.Name,
		&resource.Description,
		&resource.Capacity,
		&ownerID,
		&resource.CreatedAt,
		&resource.UpdatedAt,
		&ownerName,
	)

	if err == sql.ErrNoRows {
		return nil, ErrResourceNotFound
	}
	if err != nil {
		return nil, err
	}

	resource.OwnerID = models.NullInt64ToPtr(ownerID)
	if ownerName.Valid {
		resource.OwnerName = ownerName.String
	}

	return resource, nil
}

func (r *resourceRepository) Update(ctx context.Context, resource *models.Resource) error {
	query := `
		UPDATE resources
		SET name = $1, description = $2, capacity = $3, owner_id = $4, updated_at = $5
		WHERE id = $6
	`

	resource.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		resource.Name,
		resource.Description,
		resource.Capacity,
		resource.OwnerID,
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
		SELECT r.id, r.name, r.description, r.capacity, r.owner_id, r.created_at, r.updated_at, u.name as owner_name
		FROM resources r
		LEFT JOIN users u ON r.owner_id = u.id
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resources := make([]*models.Resource, 0)
	for rows.Next() {
		resource := &models.Resource{}
		var ownerID sql.NullInt64
		var ownerName sql.NullString
		err := rows.Scan(
			&resource.ID,
			&resource.Name,
			&resource.Description,
			&resource.Capacity,
			&ownerID,
			&resource.CreatedAt,
			&resource.UpdatedAt,
			&ownerName,
		)
		if err != nil {
			return nil, err
		}
		resource.OwnerID = models.NullInt64ToPtr(ownerID)
		if ownerName.Valid {
			resource.OwnerName = ownerName.String
		}
		resources = append(resources, resource)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}
