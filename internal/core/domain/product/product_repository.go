package product

import (
	"context"

	"github.com/ljj/gugu-api/internal/core/enum"
)

type Repository interface {
	FindByID(ctx context.Context, productID string) (*Product, error)
	FindByIDs(ctx context.Context, productIDs []string) ([]Product, error)
	FindByMarketAndExternalProductID(ctx context.Context, market enum.Market, externalProductID string) (*Product, error)
	ListByMarket(ctx context.Context, market enum.Market) ([]Product, error)
	ListByCollectionSource(ctx context.Context, source string, limit int, offset int) ([]Product, error)
	Create(ctx context.Context, product Product) error
	Update(ctx context.Context, product Product) error
}
