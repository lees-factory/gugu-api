package pricehistory

import "context"

type Finder interface {
	ListByProductID(ctx context.Context, productID string) ([]PriceHistory, error)
}

type finder struct {
	repository Repository
}

func NewFinder(repository Repository) Finder {
	return &finder{repository: repository}
}

func (f *finder) ListByProductID(ctx context.Context, productID string) ([]PriceHistory, error) {
	return f.repository.ListByProductID(ctx, productID)
}
