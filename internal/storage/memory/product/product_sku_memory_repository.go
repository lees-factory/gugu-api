package product

import (
	"context"
	"sync"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type ProductSKUMemoryRepository struct {
	mu          sync.RWMutex
	byID        map[string]domainproduct.ProductSKU
	byProductID map[string][]string
}

func NewSKURepository() *ProductSKUMemoryRepository {
	return &ProductSKUMemoryRepository{
		byID:        make(map[string]domainproduct.ProductSKU),
		byProductID: make(map[string][]string),
	}
}

func (r *ProductSKUMemoryRepository) Create(_ context.Context, sku domainproduct.ProductSKU) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[sku.ID] = sku
	r.byProductID[sku.ProductID] = append(r.byProductID[sku.ProductID], sku.ID)
	return nil
}

func (r *ProductSKUMemoryRepository) FindByID(_ context.Context, skuID string) (*domainproduct.ProductSKU, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.byID[skuID]
	if !ok {
		return nil, nil
	}
	found := item
	return &found, nil
}

func (r *ProductSKUMemoryRepository) FindByProductID(_ context.Context, productID string) ([]domainproduct.ProductSKU, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byProductID[productID]
	skus := make([]domainproduct.ProductSKU, 0, len(ids))
	for _, id := range ids {
		if sku, ok := r.byID[id]; ok {
			skus = append(skus, sku)
		}
	}
	return skus, nil
}
