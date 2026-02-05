package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"smartbooking/internal/models"
)

var (
	ErrBookingNotFound = errors.New("booking not found")
)

// BookingRepository defines the interface for booking data operations
type BookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) error
	GetByID(ctx context.Context, id int64) (*models.Booking, error)
	Update(ctx context.Context, booking *models.Booking) error
	Delete(ctx context.Context, id int64) error
	ListByUser(ctx context.Context, userID int64) ([]*models.Booking, error)
	ListByResource(ctx context.Context, resourceID int64) ([]*models.Booking, error)
	ListAll(ctx context.Context) ([]*models.Booking, error)
	CheckOverlap(ctx context.Context, resourceID int64, startTime, endTime time.Time) (bool, error)
}

// bookingRepository implements BookingRepository interface with PostgreSQL storage
type bookingRepository struct {
	db *sql.DB
}

// NewBookingRepository creates a new instance of BookingRepository
func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{
		db: db,
	}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	query := `
		INSERT INTO bookings (user_id, resource_id, start_time, end_time, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	booking.CreatedAt = now
	booking.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		booking.UserID,
		booking.ResourceID,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		booking.CreatedAt,
		booking.UpdatedAt,
	).Scan(&booking.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *bookingRepository) GetByID(ctx context.Context, id int64) (*models.Booking, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`

	booking := &models.Booking{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.ResourceID,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrBookingNotFound
	}
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (r *bookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	query := `
		UPDATE bookings
		SET user_id = $1, resource_id = $2, start_time = $3, end_time = $4, status = $5, updated_at = $6
		WHERE id = $7
	`

	booking.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		booking.UserID,
		booking.ResourceID,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		booking.UpdatedAt,
		booking.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrBookingNotFound
	}

	return nil
}

func (r *bookingRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM bookings WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrBookingNotFound
	}

	return nil
}

func (r *bookingRepository) ListByUser(ctx context.Context, userID int64) ([]*models.Booking, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]*models.Booking, 0)
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.ResourceID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) ListByResource(ctx context.Context, resourceID int64) ([]*models.Booking, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE resource_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, resourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]*models.Booking, 0)
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.ResourceID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) ListAll(ctx context.Context) ([]*models.Booking, error) {
	query := `
		SELECT id, user_id, resource_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]*models.Booking, 0)
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.ResourceID,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *bookingRepository) CheckOverlap(ctx context.Context, resourceID int64, startTime, endTime time.Time) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM bookings
		WHERE resource_id = $1
			AND status != 'cancelled'
			AND start_time < $3
			AND end_time > $2
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, resourceID, startTime, endTime).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
