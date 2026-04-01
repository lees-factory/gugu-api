package pricealert

import "context"

type Repository interface {
	FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*PriceAlert, error)
	ListByProductID(ctx context.Context, productID string) ([]PriceAlert, error)
	ListByProductIDs(ctx context.Context, productIDs []string) ([]PriceAlert, error)
	ListByUserID(ctx context.Context, userID string) ([]PriceAlert, error)
	Create(ctx context.Context, alert PriceAlert) error
	UpdateEnabled(ctx context.Context, alertID string, enabled bool) error
	DeleteByUserIDAndProductID(ctx context.Context, userID string, productID string) error
}
