package auth

import "context"

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
