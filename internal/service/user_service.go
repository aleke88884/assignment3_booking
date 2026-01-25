package service

import (
	"context"

	"smartbooking/internal/models"
	"smartbooking/internal/repository"
)

// UserService handles user-related business logic
type UserService interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) List(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.List(ctx)
}
