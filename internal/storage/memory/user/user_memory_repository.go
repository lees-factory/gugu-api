package user

import (
	"context"
	"errors"
	"sync"
	"time"

	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
)

type UserMemoryRepository struct {
	mu      sync.RWMutex
	users   map[string]domainuser.User
	userIDs map[string]string
}

func NewRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		users:   make(map[string]domainuser.User),
		userIDs: make(map[string]string),
	}
}

func (r *UserMemoryRepository) FindByEmail(_ context.Context, email string) (*domainuser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.userIDs[email]
	if !ok {
		return nil, nil
	}

	return new(r.users[id]), nil
}

func (r *UserMemoryRepository) FindByID(_ context.Context, userID string) (*domainuser.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	foundUser, ok := r.users[userID]
	if !ok {
		return nil, nil
	}

	return &foundUser, nil
}

func (r *UserMemoryRepository) Create(_ context.Context, newUser domainuser.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.userIDs[newUser.Email]; exists {
		return errors.New("user already exists")
	}

	r.users[newUser.ID] = newUser
	r.userIDs[newUser.Email] = newUser.ID
	return nil
}

func (r *UserMemoryRepository) MarkEmailVerified(_ context.Context, userID string, verifiedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	foundUser, ok := r.users[userID]
	if !ok {
		return errors.New("user not found")
	}

	foundUser.EmailVerified = true
	foundUser.EmailVerifiedAt = &verifiedAt
	r.users[userID] = foundUser
	return nil
}
