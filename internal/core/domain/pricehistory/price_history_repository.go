package pricehistory

import "context"

type Repository interface {
	Create(ctx context.Context, history PriceHistory) error
	ListByProductID(ctx context.Context, productID string, currency string) ([]PriceHistory, error)
}
