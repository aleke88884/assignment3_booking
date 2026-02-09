package repository

import (
	"context"
	"database/sql"
	"fmt"

	"smartbooking/internal/models"
)

// PhotoRepository интерфейс для работы с фотографиями
type PhotoRepository interface {
	Create(ctx context.Context, photo *models.ResourcePhoto) error
	GetByID(ctx context.Context, id int64) (*models.ResourcePhoto, error)
	GetByResourceID(ctx context.Context, resourceID int64) ([]*models.ResourcePhoto, error)
	Delete(ctx context.Context, id int64) error
	SetPrimary(ctx context.Context, id int64, resourceID int64) error
	UpdateOrder(ctx context.Context, id int64, order int) error
}

type photoRepository struct {
	db *sql.DB
}

func NewPhotoRepository(db *sql.DB) PhotoRepository {
	return &photoRepository{db: db}
}

func (r *photoRepository) Create(ctx context.Context, photo *models.ResourcePhoto) error {
	query := `
		INSERT INTO resource_photos (
			resource_id, url, storage_key, file_name, file_size,
			mime_type, width, height, is_primary, display_order
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		photo.ResourceID, photo.URL, photo.StorageKey, photo.FileName, photo.FileSize,
		photo.MimeType, photo.Width, photo.Height, photo.IsPrimary, photo.DisplayOrder,
	).Scan(&photo.ID, &photo.CreatedAt, &photo.UpdatedAt)

	if err != nil {
		return fmt.Errorf("ошибка создания фото: %w", err)
	}

	return nil
}

func (r *photoRepository) GetByID(ctx context.Context, id int64) (*models.ResourcePhoto, error) {
	query := `
		SELECT id, resource_id, url, storage_key, file_name, file_size,
			mime_type, width, height, is_primary, display_order,
			created_at, updated_at
		FROM resource_photos
		WHERE id = $1
	`

	photo := &models.ResourcePhoto{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&photo.ID, &photo.ResourceID, &photo.URL, &photo.StorageKey, &photo.FileName,
		&photo.FileSize, &photo.MimeType, &photo.Width, &photo.Height,
		&photo.IsPrimary, &photo.DisplayOrder, &photo.CreatedAt, &photo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("фото не найдено")
	}
	if err != nil {
		return nil, err
	}

	return photo, nil
}

func (r *photoRepository) GetByResourceID(ctx context.Context, resourceID int64) ([]*models.ResourcePhoto, error) {
	query := `
		SELECT id, resource_id, url, storage_key, file_name, file_size,
			mime_type, width, height, is_primary, display_order,
			created_at, updated_at
		FROM resource_photos
		WHERE resource_id = $1
		ORDER BY is_primary DESC, display_order ASC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := make([]*models.ResourcePhoto, 0)
	for rows.Next() {
		photo := &models.ResourcePhoto{}
		err := rows.Scan(
			&photo.ID, &photo.ResourceID, &photo.URL, &photo.StorageKey, &photo.FileName,
			&photo.FileSize, &photo.MimeType, &photo.Width, &photo.Height,
			&photo.IsPrimary, &photo.DisplayOrder, &photo.CreatedAt, &photo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, rows.Err()
}

func (r *photoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM resource_photos WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("фото не найдено")
	}

	return nil
}

func (r *photoRepository) SetPrimary(ctx context.Context, id int64, resourceID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Убираем is_primary у всех фото ресурса
	_, err = tx.ExecContext(ctx,
		"UPDATE resource_photos SET is_primary = false WHERE resource_id = $1",
		resourceID,
	)
	if err != nil {
		return err
	}

	// Устанавливаем is_primary для выбранного фото
	_, err = tx.ExecContext(ctx,
		"UPDATE resource_photos SET is_primary = true WHERE id = $1 AND resource_id = $2",
		id, resourceID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *photoRepository) UpdateOrder(ctx context.Context, id int64, order int) error {
	query := `UPDATE resource_photos SET display_order = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, order, id)
	return err
}
