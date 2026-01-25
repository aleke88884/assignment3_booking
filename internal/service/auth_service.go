package service

import (
	"context"
	"errors"

	"smartbooking/internal/models"
	"smartbooking/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

// AuthService handles authentication and authorization
type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	// TODO: Check if user exists
	// TODO: Hash password
	// TODO: Create user
	return nil, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*models.User, error) {
	// TODO: Get user by email
	// TODO: Verify password
	// TODO: Return user or error
	return nil, nil
}
