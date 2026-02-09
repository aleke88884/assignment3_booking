package service

import (
	"context"
	"errors"

	"smartbooking/internal/models"
	"smartbooking/internal/repository"
)

var (
	ErrInvalidRating = errors.New("rating must be between 1 and 5")
)

type ReviewService interface {
	Create(ctx context.Context, userID, resourceID int64, bookingID *int64, rating int, comment string) (*models.Review, error)
	GetByID(ctx context.Context, id int64) (*models.Review, error)
	GetByResource(ctx context.Context, resourceID int64) ([]*models.Review, error)
	GetByUser(ctx context.Context, userID int64) ([]*models.Review, error)
	Update(ctx context.Context, id int64, rating int, comment string) (*models.Review, error)
	Delete(ctx context.Context, id int64) error
	GetResourceAverageRating(ctx context.Context, resourceID int64) (float64, error)
}

type reviewService struct {
	reviewRepo reviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
	}
}

func (s *reviewService) Create(ctx context.Context, userID, resourceID int64, bookingID *int64, rating int, comment string) (*models.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}

	review := &models.Review{
		UserID:     userID,
		ResourceID: resourceID,
		BookingID:  bookingID,
		Rating:     rating,
		Comment:    comment,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) GetByID(ctx context.Context, id int64) (*models.Review, error) {
	return s.reviewRepo.GetByID(ctx, id)
}

func (s *reviewService) GetByResource(ctx context.Context, resourceID int64) ([]*models.Review, error) {
	return s.reviewRepo.GetByResource(ctx, resourceID)
}

func (s *reviewService) GetByUser(ctx context.Context, userID int64) ([]*models.Review, error) {
	return s.reviewRepo.GetByUser(ctx, userID)
}

func (s *reviewService) Update(ctx context.Context, id int64, rating int, comment string) (*models.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}

	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	review.Rating = rating
	review.Comment = comment

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) Delete(ctx context.Context, id int64) error {
	return s.reviewRepo.Delete(ctx, id)
}

func (s *reviewService) GetResourceAverageRating(ctx context.Context, resourceID int64) (float64, error) {
	return s.reviewRepo.GetAverageRating(ctx, resourceID)
}

type reviewRepository interface {
	Create(ctx context.Context, review *models.Review) error
	GetByID(ctx context.Context, id int64) (*models.Review, error)
	GetByResource(ctx context.Context, resourceID int64) ([]*models.Review, error)
	GetByUser(ctx context.Context, userID int64) ([]*models.Review, error)
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, id int64) error
	GetAverageRating(ctx context.Context, resourceID int64) (float64, error)
}
