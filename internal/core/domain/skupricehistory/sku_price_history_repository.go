package skupricehistory

import "context"

type Repository interface {
	Create(ctx context.Context, history SKUPriceHistory) error
	ListBySKUID(ctx context.Context, skuID string) ([]SKUPriceHistory, error)
}
