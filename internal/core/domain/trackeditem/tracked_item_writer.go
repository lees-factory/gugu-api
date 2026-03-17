package trackeditem

import "context"

type Writer interface {
	Create(ctx context.Context, trackedItem TrackedItem) error
	DeleteByIDAndUserID(ctx context.Context, trackedItemID string, userID string) error
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
