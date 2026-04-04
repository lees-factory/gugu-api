package product

import (
	"context"
	"strings"
	"sync"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type ProductVariantMemoryRepository struct {
	mu    sync.RWMutex
	byKey map[string]domainproduct.Variant
}

func NewVariantRepository() *ProductVariantMemoryRepository {
	return &ProductVariantMemoryRepository{
		byKey: make(map[string]domainproduct.Variant),
	}
}

func (r *ProductVariantMemoryRepository) FindByProductIDLanguageCurrency(_ context.Context, productID string, language string, currency string) (*domainproduct.Variant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.byKey[variantKey(productID, language, currency)]
	if !ok {
		return nil, nil
	}
	found := item
	return &found, nil
}

func (r *ProductVariantMemoryRepository) FindByLookupKeys(_ context.Context, keys []domainproduct.VariantLookupKey) ([]domainproduct.Variant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domainproduct.Variant, 0, len(keys))
	for _, key := range keys {
		item, ok := r.byKey[variantKey(key.ProductID, key.Language, key.Currency)]
		if !ok {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ProductVariantMemoryRepository) Upsert(_ context.Context, variant domainproduct.Variant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byKey[variantKey(variant.ProductID, variant.Language, variant.Currency)] = variant
	return nil
}

func variantKey(productID string, language string, currency string) string {
	return strings.TrimSpace(productID) + ":" + strings.ToUpper(strings.TrimSpace(language)) + ":" + strings.ToUpper(strings.TrimSpace(currency))
}
