package service

import (
	"context"

	"smartbooking/internal/models"
	"smartbooking/internal/repository"
)

// ResourceService handles resource-related business logic
type ResourceService interface {
	Create(ctx context.Context, resource *models.Resource) error
	GetByID(ctx context.Context, id int64) (*models.Resource, error)
	Update(ctx context.Context, resource *models.Resource) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.Resource, error)
}

type resourceService struct {
	resourceRepo repository.ResourceRepository
}

// NewResourceService creates a new ResourceService instance
func NewResourceService(resourceRepo repository.ResourceRepository) ResourceService {
	return &resourceService{
		resourceRepo: resourceRepo,
	}
}

func (s *resourceService) Create(ctx context.Context, resource *models.Resource) error {
	return s.resourceRepo.Create(ctx, resource)
}

func (s *resourceService) GetByID(ctx context.Context, id int64) (*models.Resource, error) {
	return s.resourceRepo.GetByID(ctx, id)
}

func (s *resourceService) Update(ctx context.Context, resource *models.Resource) error {
	return s.resourceRepo.Update(ctx, resource)
}

func (s *resourceService) Delete(ctx context.Context, id int64) error {
	return s.resourceRepo.Delete(ctx, id)
}

func (s *resourceService) List(ctx context.Context) ([]*models.Resource, error) {
	return s.resourceRepo.List(ctx)
}
