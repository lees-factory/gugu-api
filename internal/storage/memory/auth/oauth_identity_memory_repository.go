package auth

import (
	"context"
	"errors"
	"sync"
	"time"

	supportauth "github.com/ljj/gugu-api/internal/support/auth"
)

type OAuthIdentityMemoryRepository struct {
	mu         sync.RWMutex
	identities map[string]supportauth.OAuthIdentity
}

func NewOAuthIdentityRepository() *OAuthIdentityMemoryRepository {
	return &OAuthIdentityMemoryRepository{
		identities: make(map[string]supportauth.OAuthIdentity),
	}
}

func (r *OAuthIdentityMemoryRepository) FindByProviderSubject(_ context.Context, provider string, subject string) (*supportauth.OAuthIdentity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	foundIdentity, ok := r.identities[provider+":"+subject]
	if !ok {
		return nil, nil
	}

	return &foundIdentity, nil
}

func (r *OAuthIdentityMemoryRepository) Create(_ context.Context, identity supportauth.OAuthIdentity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.identities[identity.Provider+":"+identity.Subject] = identity
	return nil
}

func (r *OAuthIdentityMemoryRepository) UpdateLastLogin(_ context.Context, provider string, subject string, lastLoginAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := provider + ":" + subject
	foundIdentity, ok := r.identities[key]
	if !ok {
		return errors.New("oauth identity not found")
	}

	foundIdentity.LastLoginAt = lastLoginAt
	r.identities[key] = foundIdentity
	return nil
}
