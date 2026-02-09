package service

import (
	"context"
	"smartbooking/internal/models"
	"smartbooking/internal/repository"
)

// OwnerService provides business logic for owner operations
type OwnerService interface {
	GetOwnerResources(ctx context.Context, ownerID int64) ([]*models.Resource, error)
	GetOwnerBookings(ctx context.Context, ownerID int64) ([]*models.Booking, error)
	GetOwnerStatistics(ctx context.Context, ownerID int64) (*repository.OwnerStatistics, error)
}

type ownerService struct {
	ownerRepo repository.OwnerRepository
}

// NewOwnerService creates a new OwnerService
func NewOwnerService(ownerRepo repository.OwnerRepository) OwnerService {
	return &ownerService{
		ownerRepo: ownerRepo,
	}
}

// GetOwnerResources retrieves all resources for an owner
func (s *ownerService) GetOwnerResources(ctx context.Context, ownerID int64) ([]*models.Resource, error) {
	return s.ownerRepo.GetOwnerResources(ctx, ownerID)
}

// GetOwnerBookings retrieves all bookings for an owner's resources
func (s *ownerService) GetOwnerBookings(ctx context.Context, ownerID int64) ([]*models.Booking, error) {
	return s.ownerRepo.GetOwnerBookings(ctx, ownerID)
}

// GetOwnerStatistics retrieves statistics for an owner
func (s *ownerService) GetOwnerStatistics(ctx context.Context, ownerID int64) (*repository.OwnerStatistics, error) {
	return s.ownerRepo.GetOwnerStatistics(ctx, ownerID)
}
