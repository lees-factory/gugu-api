package auth

import (
	"context"
	"errors"
	"sync"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

type OAuthIdentityMemoryRepository struct {
	mu         sync.RWMutex
	identities map[string]domainauth.OAuthIdentity
}

func NewOAuthIdentityRepository() *OAuthIdentityMemoryRepository {
	return &OAuthIdentityMemoryRepository{
		identities: make(map[string]domainauth.OAuthIdentity),
	}
}

func (r *OAuthIdentityMemoryRepository) FindByProviderSubject(_ context.Context, provider string, subject string) (*domainauth.OAuthIdentity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	foundIdentity, ok := r.identities[provider+":"+subject]
	if !ok {
		return nil, nil
	}

	return &foundIdentity, nil
}

func (r *OAuthIdentityMemoryRepository) Create(_ context.Context, identity domainauth.OAuthIdentity) error {
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
