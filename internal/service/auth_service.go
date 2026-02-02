package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      models.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*models.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
