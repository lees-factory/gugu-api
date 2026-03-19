package product

import "context"

type Repository interface {
	FindByID(ctx context.Context, productID string) (*Product, error)
	FindByIDs(ctx context.Context, productIDs []string) ([]Product, error)
	FindByMarketAndExternalProductID(ctx context.Context, market Market, externalProductID string) (*Product, error)
	Create(ctx context.Context, product Product) error
	Update(ctx context.Context, product Product) error
}
