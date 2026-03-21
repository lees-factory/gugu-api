package pricehistory

import "context"

type Writer interface {
	Create(ctx context.Context, history PriceHistory) error
}

type writer struct {
	repository Repository
}

func NewWriter(repository Repository) Writer {
	return &writer{repository: repository}
}

func (w *writer) Create(ctx context.Context, history PriceHistory) error {
	return w.repository.Create(ctx, history)
}
