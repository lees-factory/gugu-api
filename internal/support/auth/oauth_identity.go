package auth

import (
	"context"
	"time"
)

type OAuthProvider string

const (
	OAuthProviderGoogle OAuthProvider = "google"
)

type OAuthIdentity struct {
	ID          string
	UserID      string
	Provider    string
	Subject     string
	Email       string
	CreatedAt   time.Time
	LastLoginAt time.Time
}

type OAuthIdentityRepository interface {
	FindByProviderSubject(ctx context.Context, provider string, subject string) (*OAuthIdentity, error)
	Create(ctx context.Context, identity OAuthIdentity) error
	UpdateLastLogin(ctx context.Context, provider string, subject string, lastLoginAt time.Time) error
}

type OAuthIdentityFinder interface {
	FindByProviderSubject(ctx context.Context, provider string, subject string) (*OAuthIdentity, error)
}

type oauthIdentityFinder struct {
	repository OAuthIdentityRepository
}

func NewOAuthIdentityFinder(repository OAuthIdentityRepository) OAuthIdentityFinder {
	return &oauthIdentityFinder{repository: repository}
}

func (f *oauthIdentityFinder) FindByProviderSubject(ctx context.Context, provider, subject string) (*OAuthIdentity, error) {
	return f.repository.FindByProviderSubject(ctx, provider, subject)
}

type OAuthIdentityWriter interface {
	Create(ctx context.Context, identity OAuthIdentity) error
	UpdateLastLogin(ctx context.Context, provider string, subject string, lastLoginAt time.Time) error
}

type oauthIdentityWriter struct {
	repository OAuthIdentityRepository
}

func NewOAuthIdentityWriter(repository OAuthIdentityRepository) OAuthIdentityWriter {
	return &oauthIdentityWriter{repository: repository}
}

func (w *oauthIdentityWriter) Create(ctx context.Context, identity OAuthIdentity) error {
	return w.repository.Create(ctx, identity)
}

func (w *oauthIdentityWriter) UpdateLastLogin(ctx context.Context, provider string, subject string, lastLoginAt time.Time) error {
	return w.repository.UpdateLastLogin(ctx, provider, subject, lastLoginAt)
}
