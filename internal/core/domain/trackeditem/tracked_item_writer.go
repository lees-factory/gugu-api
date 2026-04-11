package trackeditem

import "context"

type Writer interface {
	Create(ctx context.Context, trackedItem TrackedItem) error
	DeleteByIDAndUserID(ctx context.Context, trackedItemID string, userID string) error
	UpdateSKU(ctx context.Context, trackedItemID string, userID string, skuID string) error
	UpdatePreferredLanguage(ctx context.Context, trackedItemID string, userID string, preferredLanguage string) error
	UpdateTrackingScope(ctx context.Context, trackedItemID string, userID string, trackingScope string) error
	ReplaceWatchSKUs(ctx context.Context, trackedItemID string, skuIDs []string) error
}

type writer struct {
	repository Repository
}

func NewWriter(repository Repository) Writer {
	return &writer{repository: repository}
}

func (w *writer) Create(ctx context.Context, trackedItem TrackedItem) error {
	return w.repository.Create(ctx, trackedItem)
}

func (w *writer) DeleteByIDAndUserID(ctx context.Context, trackedItemID string, userID string) error {
	return w.repository.DeleteByIDAndUserID(ctx, trackedItemID, userID)
}

func (w *writer) UpdateSKU(ctx context.Context, trackedItemID string, userID string, skuID string) error {
	return w.repository.UpdateSKU(ctx, trackedItemID, userID, skuID)
}

func (w *writer) UpdatePreferredLanguage(ctx context.Context, trackedItemID string, userID string, preferredLanguage string) error {
	return w.repository.UpdatePreferredLanguage(ctx, trackedItemID, userID, preferredLanguage)
}

func (w *writer) UpdateTrackingScope(ctx context.Context, trackedItemID string, userID string, trackingScope string) error {
	return w.repository.UpdateTrackingScope(ctx, trackedItemID, userID, trackingScope)
}

func (w *writer) ReplaceWatchSKUs(ctx context.Context, trackedItemID string, skuIDs []string) error {
	return w.repository.ReplaceWatchSKUs(ctx, trackedItemID, skuIDs)
}
