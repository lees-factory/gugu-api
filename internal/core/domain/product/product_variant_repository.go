package product

import "context"

type VariantLookupKey struct {
	ProductID string
	Language  string
	Currency  string
}

type VariantRepository interface {
	FindByProductIDLanguageCurrency(ctx context.Context, productID string, language string, currency string) (*Variant, error)
	FindByLookupKeys(ctx context.Context, keys []VariantLookupKey) ([]Variant, error)
	Upsert(ctx context.Context, variant Variant) error
}
