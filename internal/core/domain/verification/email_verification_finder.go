package verification

import "context"

type Finder interface {
	FindByToken(ctx context.Context, token string) (*EmailVerification, error)
}

type finder struct {
	repository Repository
}

func NewFinder(repository Repository) Finder {
	return &finder{repository: repository}
}

func (f *finder) FindByToken(ctx context.Context, token string) (*EmailVerification, error) {
	return f.repository.FindByToken(ctx, token)
}
