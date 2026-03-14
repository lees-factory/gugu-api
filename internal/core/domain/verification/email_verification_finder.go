package verification

import "context"

type Finder interface {
	FindByCode(ctx context.Context, code string) (*EmailVerification, error)
}

type finder struct {
	repository Repository
}

func NewFinder(repository Repository) Finder {
	return &finder{repository: repository}
}

func (f *finder) FindByCode(ctx context.Context, code string) (*EmailVerification, error) {
	return f.repository.FindByCode(ctx, code)
}
