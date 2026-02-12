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

	// Load photos for all resources
	if len(resources) > 0 {
		err = r.loadPhotosForResources(ctx, resources)
		if err != nil {
			// Don't fail if photos can't be loaded, just log
			// return nil, err
		}
	}

	return resources, nil
}

func (r *resourceRepository) loadPhotosForResources(ctx context.Context, resources []*models.Resource) error {
	if len(resources) == 0 {
		return nil
	}

	// Get all resource IDs
	resourceIDs := make([]interface{}, len(resources))
	resourceMap := make(map[int64]*models.Resource)
	for i, res := range resources {
		resourceIDs[i] = res.ID
		resourceMap[res.ID] = res
		res.Photos = make([]models.ResourcePhoto, 0)
	}

	// Build query with IN clause
	query := `
		SELECT id, resource_id, url, storage_key, file_name, file_size, mime_type,
		       width, height, is_primary, display_order, created_at, updated_at
		FROM resource_photos
		WHERE resource_id = ANY($1)
		ORDER BY resource_id, display_order, id
	`

	rows, err := r.db.QueryContext(ctx, query, resourceIDs)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		photo := models.ResourcePhoto{}
		var width, height sql.NullInt64
		err := rows.Scan(
			&photo.ID,
			&photo.ResourceID,
			&photo.URL,
			&photo.StorageKey,
			&photo.FileName,
			&photo.FileSize,
			&photo.MimeType,
			&width,
			&height,
			&photo.IsPrimary,
			&photo.DisplayOrder,
			&photo.CreatedAt,
			&photo.UpdatedAt,
		)
		if err != nil {
			return err
		}

		if width.Valid {
			photo.Width = int(width.Int64)
		}
		if height.Valid {
			photo.Height = int(height.Int64)
		}

		if resource, ok := resourceMap[photo.ResourceID]; ok {
			resource.Photos = append(resource.Photos, photo)
		}
	}

	return rows.Err()
}
