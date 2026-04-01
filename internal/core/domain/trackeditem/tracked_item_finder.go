package trackeditem

import (
	"context"
	"time"
)

type Finder interface {
	FindByIDAndUserID(ctx context.Context, trackedItemID string, userID string) (*TrackedItem, error)
	FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*TrackedItem, error)
	ListByUserID(ctx context.Context, userID string) ([]TrackedItem, error)
	ListByUserIDWithCursor(ctx context.Context, userID string, cursorCreatedAt time.Time, cursorID string, limit int) ([]TrackedItem, error)
	ListByUserIDFirstPage(ctx context.Context, userID string, limit int) ([]TrackedItem, error)
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

func (f *finder) ListByUserIDWithCursor(ctx context.Context, userID string, cursorCreatedAt time.Time, cursorID string, limit int) ([]TrackedItem, error) {
	return f.repository.ListByUserIDWithCursor(ctx, userID, cursorCreatedAt, cursorID, limit)
}

func (f *finder) ListByUserIDFirstPage(ctx context.Context, userID string, limit int) ([]TrackedItem, error) {
	return f.repository.ListByUserIDFirstPage(ctx, userID, limit)
}
