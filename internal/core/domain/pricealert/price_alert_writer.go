package pricealert

import "context"

type Writer interface {
	Create(ctx context.Context, alert PriceAlert) error
	UpdateEnabled(ctx context.Context, alertID string, enabled bool) error
	UpdateSettings(ctx context.Context, alertID string, channel string, enabled bool) error
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

func (w *writer) UpdateSettings(ctx context.Context, alertID string, channel string, enabled bool) error {
	return w.repository.UpdateSettings(ctx, alertID, channel, enabled)
}
