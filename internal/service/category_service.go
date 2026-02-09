package service

import (
	"context"

	"smartbooking/internal/models"
	"smartbooking/internal/repository"
)

type CategoryService interface {
	Create(ctx context.Context, name, slug, description, icon string) (*models.ResourceCategory, error)
	GetByID(ctx context.Context, id int64) (*models.ResourceCategory, error)
	List(ctx context.Context) ([]*models.ResourceCategory, error)
	Update(ctx context.Context, id int64, name, slug, description, icon string) (*models.ResourceCategory, error)
	Delete(ctx context.Context, id int64) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) Create(ctx context.Context, name, slug, description, icon string) (*models.ResourceCategory, error) {
	category := &models.ResourceCategory{
		Name:        name,
		Slug:        slug,
		Description: description,
		Icon:        icon,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetByID(ctx context.Context, id int64) (*models.ResourceCategory, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *categoryService) List(ctx context.Context) ([]*models.ResourceCategory, error) {
	return s.categoryRepo.List(ctx)
}

func (s *categoryService) Update(ctx context.Context, id int64, name, slug, description, icon string) (*models.ResourceCategory, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	category.Name = name
	category.Slug = slug
	category.Description = description
	category.Icon = icon

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	return s.categoryRepo.Delete(ctx, id)
}
