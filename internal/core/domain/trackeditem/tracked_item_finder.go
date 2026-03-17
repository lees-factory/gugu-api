package trackeditem

import "context"

type Finder interface {
	FindByIDAndUserID(ctx context.Context, trackedItemID string, userID string) (*TrackedItem, error)
	FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*TrackedItem, error)
	ListByUserID(ctx context.Context, userID string) ([]TrackedItem, error)
}

type finder struct {
	repository Repository
}

func NewFinder(repository Repository) Finder {
	return &finder{repository: repository}
}

func (f *finder) FindByIDAndUserID(ctx context.Context, trackedItemID string, userID string) (*TrackedItem, error) {
	return f.repository.FindByIDAndUserID(ctx, trackedItemID, userID)
}

func (f *finder) FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*TrackedItem, error) {
	return f.repository.FindByUserIDAndProductID(ctx, userID, productID)
}

func (f *finder) ListByUserID(ctx context.Context, userID string) ([]TrackedItem, error) {
	return f.repository.ListByUserID(ctx, userID)
}
