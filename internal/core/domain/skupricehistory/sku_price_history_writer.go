package skupricehistory

import "context"

type Writer interface {
	Create(ctx context.Context, history SKUPriceHistory) error
}

type writer struct {
	repository Repository
}

func NewWriter(repository Repository) Writer {
	return &writer{repository: repository}
}

func (w *writer) Create(ctx context.Context, history SKUPriceHistory) error {
	return w.repository.Create(ctx, history)
}
