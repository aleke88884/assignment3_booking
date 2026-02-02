package repository

import (
	"context"
	"errors"
	"sync"
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

// bookingRepository implements BookingRepository interface with in-memory storage
type bookingRepository struct {
	mu       sync.RWMutex
	bookings map[int64]*models.Booking
	nextID   int64
}

// NewBookingRepository creates a new instance of BookingRepository
func NewBookingRepository() BookingRepository {
	return &bookingRepository{
		bookings: make(map[int64]*models.Booking),
		nextID:   1,
	}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Assign ID and store
	booking.ID = r.nextID
	r.nextID++

	r.bookings[booking.ID] = booking
	return nil
}

func (r *bookingRepository) GetByID(ctx context.Context, id int64) (*models.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	booking, exists := r.bookings[id]
	if !exists {
		return nil, ErrBookingNotFound
	}

	return booking, nil
}

func (r *bookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bookings[booking.ID]; !exists {
		return ErrBookingNotFound
	}

	r.bookings[booking.ID] = booking
	return nil
}

func (r *bookingRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bookings[id]; !exists {
		return ErrBookingNotFound
	}

	delete(r.bookings, id)
	return nil
}

func (r *bookingRepository) ListByUser(ctx context.Context, userID int64) ([]*models.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bookings := make([]*models.Booking, 0)
	for _, booking := range r.bookings {
		if booking.UserID == userID {
			bookings = append(bookings, booking)
		}
	}

	return bookings, nil
}

func (r *bookingRepository) ListByResource(ctx context.Context, resourceID int64) ([]*models.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bookings := make([]*models.Booking, 0)
	for _, booking := range r.bookings {
		if booking.ResourceID == resourceID {
			bookings = append(bookings, booking)
		}
	}

	return bookings, nil
}

func (r *bookingRepository) ListAll(ctx context.Context) ([]*models.Booking, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bookings := make([]*models.Booking, 0, len(r.bookings))
	for _, booking := range r.bookings {
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *bookingRepository) CheckOverlap(ctx context.Context, resourceID int64, startTime, endTime time.Time) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check for overlapping bookings for the same resource
	// Bookings overlap if: (startTime < existing.EndTime) && (endTime > existing.StartTime)
	for _, booking := range r.bookings {
		if booking.ResourceID == resourceID && booking.Status != models.StatusCancelled {
			if startTime.Before(booking.EndTime) && endTime.After(booking.StartTime) {
				return true, nil
			}
		}
	}

	return false, nil
}
