package repository

import (
	"context"
	"errors"
	"sync"

	"smartbooking/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]*models.User, error)
}

// userRepository implements UserRepository interface with in-memory storage
type userRepository struct {
	mu         sync.RWMutex
	users      map[int64]*models.User
	emailIndex map[string]*models.User
	nextID     int64
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository() UserRepository {
	return &userRepository{
		users:      make(map[int64]*models.User),
		emailIndex: make(map[string]*models.User),
		nextID:     1,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	if _, exists := r.emailIndex[user.Email]; exists {
		return errors.New("user with this email already exists")
	}

	// Assign ID and store
	user.ID = r.nextID
	r.nextID++

	r.users[user.ID] = user
	r.emailIndex[user.Email] = user

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.emailIndex[email]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.users[user.ID]
	if !exists {
		return ErrUserNotFound
	}

	// Update email index if email changed
	if existing.Email != user.Email {
		delete(r.emailIndex, existing.Email)
		r.emailIndex[user.Email] = user
	}

	r.users[user.ID] = user
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	delete(r.emailIndex, user.Email)
	return nil
}

func (r *userRepository) List(ctx context.Context) ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}
