package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"smartbooking/internal/models"
)

var (
	ErrReviewNotFound = errors.New("review not found")
)

type ReviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id int64) (*models.Review, error)
	GetByResource(ctx context.Context, resourceID int64) ([]*models.Review, error)
	GetByUser(ctx context.Context, userID int64) ([]*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id int64) error
	GetAverageRating(ctx context.Context, resourceID int64) (float64, error)
}

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}

func (r *reviewRepository) Create(ctx context.Context, review *models.Review) error {
	query := `
		INSERT INTO reviews (user_id, resource_id, booking_id, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		review.UserID,
		review.ResourceID,
		review.BookingID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
		review.UpdatedAt,
	).Scan(&review.ID)

	return err
}

func (r *reviewRepository) GetByID(ctx context.Context, id int64) (*models.Review, error) {
	query := `
		SELECT id, user_id, resource_id, booking_id, rating, comment, created_at, updated_at
		FROM reviews
		WHERE id = $1
	`

	review := &models.Review{}
	var bookingID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.UserID,
		&review.ResourceID,
		&bookingID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if bookingID.Valid {
		review.BookingID = &bookingID.Int64
	}

	if err == sql.ErrNoRows {
		return nil, ErrReviewNotFound
	}
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *reviewRepository) GetByResource(ctx context.Context, resourceID int64) ([]*models.Review, error) {
	query := `
		SELECT id, user_id, resource_id, booking_id, rating, comment, created_at, updated_at
		FROM reviews
		WHERE resource_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*models.Review, 0)
	for rows.Next() {
		review := &models.Review{}
		var bookingID sql.NullInt64

		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.ResourceID,
			&bookingID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if bookingID.Valid {
			review.BookingID = &bookingID.Int64
		}

		reviews = append(reviews, review)
	}

	return reviews, rows.Err()
}

func (r *reviewRepository) GetByUser(ctx context.Context, userID int64) ([]*models.Review, error) {
	query := `
		SELECT id, user_id, resource_id, booking_id, rating, comment, created_at, updated_at
		FROM reviews
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]*models.Review, 0)
	for rows.Next() {
		review := &models.Review{}
		var bookingID sql.NullInt64

		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.ResourceID,
			&bookingID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if bookingID.Valid {
			review.BookingID = &bookingID.Int64
		}

		reviews = append(reviews, review)
	}

	return reviews, rows.Err()
}

func (r *reviewRepository) Update(ctx context.Context, review *models.Review) error {
	query := `
		UPDATE reviews
		SET rating = $1, comment = $2, updated_at = $3
		WHERE id = $4
	`

	review.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		review.Rating,
		review.Comment,
		review.UpdatedAt,
		review.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrReviewNotFound
	}

	return nil
}

func (r *reviewRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM reviews WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrReviewNotFound
	}

	return nil
}

func (r *reviewRepository) GetAverageRating(ctx context.Context, resourceID int64) (float64, error) {
	query := `
		SELECT COALESCE(AVG(rating), 0)
		FROM reviews
		WHERE resource_id = $1
	`

	var avgRating float64
	err := r.db.QueryRowContext(ctx, query, resourceID).Scan(&avgRating)
	if err != nil {
		return 0, err
	}

	return avgRating, nil
}
