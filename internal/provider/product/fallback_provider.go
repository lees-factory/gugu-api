package product

import (
	"context"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type FallbackProvider struct {
	providers []domainproduct.ProductProvider
}

func NewFallbackProvider(providers ...domainproduct.ProductProvider) *FallbackProvider {
	return &FallbackProvider{providers: providers}
}

func (p *FallbackProvider) Provide(ctx context.Context, market enum.Market, externalProductID string, originalURL string) (*domainproduct.NewProduct, error) {
	for _, provider := range p.providers {
		result, err := provider.Provide(ctx, market, externalProductID, originalURL)
		if err == nil && result != nil {
			return result, nil
		}
	}
	return nil, coreerror.New(coreerror.ProductNotFound)
}
