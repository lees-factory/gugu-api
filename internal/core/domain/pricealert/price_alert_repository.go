package pricealert

import "context"

type Repository interface {
	FindByUserIDAndSKUID(ctx context.Context, userID string, skuID string) (*PriceAlert, error)
	ListBySKUID(ctx context.Context, skuID string) ([]PriceAlert, error)
	ListByProductID(ctx context.Context, productID string) ([]PriceAlert, error)
	ListByProductIDs(ctx context.Context, productIDs []string) ([]PriceAlert, error)
	ListByUserID(ctx context.Context, userID string) ([]PriceAlert, error)
	Create(ctx context.Context, alert PriceAlert) error
	UpdateEnabled(ctx context.Context, alertID string, enabled bool) error
	UpdateSettings(ctx context.Context, alertID string, channel string, enabled bool) error
}
