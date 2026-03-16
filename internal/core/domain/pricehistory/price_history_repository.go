package pricehistory

import "context"

type Repository interface {
	ListByProductID(ctx context.Context, productID string) ([]PriceHistory, error)
}
