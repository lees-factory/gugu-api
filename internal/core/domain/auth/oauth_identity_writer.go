package auth

import (
	"context"
	"time"
)

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
