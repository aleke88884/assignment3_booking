package service

import (
	"context"
	"errors"
	"time"

	"smartbooking/internal/models"
	"smartbooking/internal/repository"
)

var (
	ErrBookingConflict  = errors.New("booking conflicts with existing reservation")
	ErrInvalidTimeRange = errors.New("invalid time range")
	ErrBookingNotFound  = errors.New("booking not found")
)

// BookingService handles booking-related business logic
type BookingService interface {
	Create(ctx context.Context, userID, resourceID int64, startTime, endTime time.Time) (*models.Booking, error)
	GetByID(ctx context.Context, id int64) (*models.Booking, error)
	Cancel(ctx context.Context, id int64) error
	ListByUser(ctx context.Context, userID int64) ([]*models.Booking, error)
	ListByResource(ctx context.Context, resourceID int64) ([]*models.Booking, error)
	ListAll(ctx context.Context) ([]*models.Booking, error)
}

type bookingService struct {
	bookingRepo  repository.BookingRepository
	resourceRepo repository.ResourceRepository
}

// NewBookingService creates a new BookingService instance
func NewBookingService(bookingRepo repository.BookingRepository, resourceRepo repository.ResourceRepository) BookingService {
	return &bookingService{
		bookingRepo:  bookingRepo,
		resourceRepo: resourceRepo,
	}
}

func (s *bookingService) Create(ctx context.Context, userID, resourceID int64, startTime, endTime time.Time) (*models.Booking, error) {
	// Validate time range
	if startTime.After(endTime) || startTime.Equal(endTime) {
		return nil, ErrInvalidTimeRange
	}

	// Check for overlapping bookings (double booking prevention)
	hasOverlap, err := s.bookingRepo.CheckOverlap(ctx, resourceID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	if hasOverlap {
		return nil, ErrBookingConflict
	}

	booking := &models.Booking{
		UserID:     userID,
		ResourceID: resourceID,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) GetByID(ctx context.Context, id int64) (*models.Booking, error) {
	return s.bookingRepo.GetByID(ctx, id)
}

func (s *bookingService) Cancel(ctx context.Context, id int64) error {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if booking == nil {
		return ErrBookingNotFound
	}

	booking.Status = models.StatusCancelled
	booking.UpdatedAt = time.Now()

	return s.bookingRepo.Update(ctx, booking)
}

func (s *bookingService) ListByUser(ctx context.Context, userID int64) ([]*models.Booking, error) {
	return s.bookingRepo.ListByUser(ctx, userID)
}

func (s *bookingService) ListByResource(ctx context.Context, resourceID int64) ([]*models.Booking, error) {
	return s.bookingRepo.ListByResource(ctx, resourceID)
}

func (s *bookingService) ListAll(ctx context.Context) ([]*models.Booking, error) {
	return s.bookingRepo.ListAll(ctx)
}
