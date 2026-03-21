package product

import (
	"context"

	"github.com/ljj/gugu-api/internal/core/enum"
)

type ProductProvider interface {
	Provide(ctx context.Context, market enum.Market, externalProductID string, originalURL string) (*NewProduct, error)
}
