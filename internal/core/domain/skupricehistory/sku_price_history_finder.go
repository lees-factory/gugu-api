package skupricehistory

import "context"

type Finder interface {
	ListBySKUID(ctx context.Context, skuID string) ([]SKUPriceHistory, error)
}

type finder struct {
	repository Repository
}

func NewFinder(repository Repository) Finder {
	return &finder{repository: repository}
}

func (f *finder) ListBySKUID(ctx context.Context, skuID string) ([]SKUPriceHistory, error) {
	return f.repository.ListBySKUID(ctx, skuID)
}
