package user

import "context"

type Finder interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, userID string) (*User, error)
}

type finder struct {
	repository Repository
}

func NewFinder(repository Repository) Finder {
	return &finder{repository: repository}
}

func (f *finder) FindByEmail(ctx context.Context, email string) (*User, error) {
	return f.repository.FindByEmail(ctx, email)
}

func (f *finder) FindByID(ctx context.Context, userID string) (*User, error) {
	return f.repository.FindByID(ctx, userID)
}
