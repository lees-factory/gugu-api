package trackeditem

import (
	"context"
	"time"
)

type Repository interface {
	FindByIDAndUserID(ctx context.Context, trackedItemID string, userID string) (*TrackedItem, error)
	FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*TrackedItem, error)
	ListByUserID(ctx context.Context, userID string) ([]TrackedItem, error)
	ListByUserIDWithCursor(ctx context.Context, userID string, cursorCreatedAt time.Time, cursorID string, limit int) ([]TrackedItem, error)
	ListByUserIDFirstPage(ctx context.Context, userID string, limit int) ([]TrackedItem, error)
	Create(ctx context.Context, trackedItem TrackedItem) error
	DeleteByIDAndUserID(ctx context.Context, trackedItemID string, userID string) error
	UpdateSKU(ctx context.Context, trackedItemID string, userID string, skuID string) error
}
