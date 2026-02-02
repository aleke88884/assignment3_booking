package repository

import (
	"context"
	"errors"
	"sync"

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

// resourceRepository implements ResourceRepository interface with in-memory storage
type resourceRepository struct {
	mu        sync.RWMutex
	resources map[int64]*models.Resource
	nextID    int64
}

// NewResourceRepository creates a new instance of ResourceRepository
func NewResourceRepository() ResourceRepository {
	return &resourceRepository{
		resources: make(map[int64]*models.Resource),
		nextID:    1,
	}
}

func (r *resourceRepository) Create(ctx context.Context, resource *models.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Assign ID and store
	resource.ID = r.nextID
	r.nextID++

	r.resources[resource.ID] = resource
	return nil
}

func (r *resourceRepository) GetByID(ctx context.Context, id int64) (*models.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resource, exists := r.resources[id]
	if !exists {
		return nil, ErrResourceNotFound
	}

	return resource, nil
}

func (r *resourceRepository) Update(ctx context.Context, resource *models.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[resource.ID]; !exists {
		return ErrResourceNotFound
	}

	r.resources[resource.ID] = resource
	return nil
}

func (r *resourceRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[id]; !exists {
		return ErrResourceNotFound
	}

	delete(r.resources, id)
	return nil
}

func (r *resourceRepository) List(ctx context.Context) ([]*models.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resources := make([]*models.Resource, 0, len(r.resources))
	for _, resource := range r.resources {
		resources = append(resources, resource)
	}

	return resources, nil
}
