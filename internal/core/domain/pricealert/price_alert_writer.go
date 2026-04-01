package pricealert

import "context"

type Writer interface {
	Create(ctx context.Context, alert PriceAlert) error
	UpdateEnabled(ctx context.Context, alertID string, enabled bool) error
	DeleteByUserIDAndProductID(ctx context.Context, userID string, productID string) error
}

type writer struct {
	repository Repository
}

func NewWriter(repository Repository) Writer {
	return &writer{repository: repository}
}

func (w *writer) Create(ctx context.Context, alert PriceAlert) error {
	return w.repository.Create(ctx, alert)
}

func (w *writer) UpdateEnabled(ctx context.Context, alertID string, enabled bool) error {
	return w.repository.UpdateEnabled(ctx, alertID, enabled)
}

func (w *writer) DeleteByUserIDAndProductID(ctx context.Context, userID string, productID string) error {
	return w.repository.DeleteByUserIDAndProductID(ctx, userID, productID)
}
