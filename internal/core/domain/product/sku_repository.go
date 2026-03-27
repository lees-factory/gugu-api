package product

import "context"

type SKURepository interface {
	Create(ctx context.Context, sku SKU) error
	FindByID(ctx context.Context, skuID string) (*SKU, error)
	FindByProductID(ctx context.Context, productID string) ([]SKU, error)
}
