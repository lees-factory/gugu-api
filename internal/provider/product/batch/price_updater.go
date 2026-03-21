package batch

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	provideraliexpress "github.com/ljj/gugu-api/internal/provider/product/aliexpress"
)

type PriceUpdater struct {
	productService     *domainproduct.Service
	priceHistoryWriter domainpricehistory.Writer
	aliExpressFetcher  *provideraliexpress.BatchFetcher
	clock              func() time.Time
}

func NewPriceUpdater(
	productService *domainproduct.Service,
	priceHistoryWriter domainpricehistory.Writer,
	aliExpressFetcher *provideraliexpress.BatchFetcher,
) *PriceUpdater {
	return &PriceUpdater{
		productService:     productService,
		priceHistoryWriter: priceHistoryWriter,
		aliExpressFetcher:  aliExpressFetcher,
		clock:              func() time.Time { return time.Now() },
	}
}

func (u *PriceUpdater) UpdateAll(ctx context.Context) error {
	if err := u.updateMarket(ctx, enum.MarketAliExpress); err != nil {
		return fmt.Errorf("update aliexpress prices: %w", err)
	}
	return nil
}

func (u *PriceUpdater) updateMarket(ctx context.Context, market enum.Market) error {
	products, err := u.productService.ListByMarket(ctx, market)
	if err != nil {
		return fmt.Errorf("list products by market: %w", err)
	}
	if len(products) == 0 {
		return nil
	}

	productMap := make(map[string]domainproduct.Product, len(products))
	externalIDs := make([]string, len(products))
	for i, p := range products {
		productMap[p.ExternalProductID] = p
		externalIDs[i] = p.ExternalProductID
	}

	var prices []provideraliexpress.PriceResult
	switch market {
	case enum.MarketAliExpress:
		prices, err = u.aliExpressFetcher.FetchPrices(ctx, externalIDs)
		if err != nil {
			return fmt.Errorf("fetch aliexpress prices: %w", err)
		}
	default:
		log.Printf("batch price update not supported for market: %s", market)
		return nil
	}

	now := u.clock()
	updated := 0
	for _, pr := range prices {
		product, ok := productMap[pr.ExternalProductID]
		if !ok {
			continue
		}
		if pr.Price == "" {
			continue
		}
		if product.CurrentPrice == pr.Price && product.Currency == pr.Currency {
			continue
		}

		changeValue := calculateChange(product.CurrentPrice, pr.Price)

		if err := u.productService.UpdatePrice(ctx, product.ID, pr.Price, pr.Currency); err != nil {
			log.Printf("failed to update price for product %s: %v", product.ID, err)
			continue
		}

		if err := u.priceHistoryWriter.Create(ctx, domainpricehistory.PriceHistory{
			ProductID:   product.ID,
			Price:       pr.Price,
			Currency:    pr.Currency,
			RecordedAt:  now,
			ChangeValue: changeValue,
		}); err != nil {
			log.Printf("failed to record price history for product %s: %v", product.ID, err)
			continue
		}

		updated++
	}

	log.Printf("batch price update: market=%s total=%d fetched=%d updated=%d", market, len(products), len(prices), updated)
	return nil
}

func calculateChange(oldPrice string, newPrice string) string {
	if oldPrice == "" || newPrice == "" {
		return "0"
	}
	old, ok1 := new(big.Float).SetString(oldPrice)
	cur, ok2 := new(big.Float).SetString(newPrice)
	if !ok1 || !ok2 {
		return "0"
	}
	diff := new(big.Float).Sub(cur, old)
	return diff.Text('f', -1)
}
