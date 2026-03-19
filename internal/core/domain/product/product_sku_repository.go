package product

import "context"

type SKURepository interface {
	Create(ctx context.Context, sku ProductSKU) error
	FindByID(ctx context.Context, skuID string) (*ProductSKU, error)
	FindByProductID(ctx context.Context, productID string) ([]ProductSKU, error)
}
