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

var supportedCurrencies = []string{"KRW", "USD"}

type PriceChangeNotifier interface {
	NotifyPriceChange(ctx context.Context, skuID string, productTitle string, oldPrice string, newPrice string, currency string)
}

type PriceFetcher interface {
	FetchPrices(ctx context.Context, externalProductIDs []string, currency string) ([]provideraliexpress.PriceResult, error)
}

type PriceUpdater struct {
	productService     *domainproduct.Service
	priceHistoryWriter domainpricehistory.Writer
	fetcher            PriceFetcher
	notifier           PriceChangeNotifier
	clock              func() time.Time
}

func NewPriceUpdater(
	productService *domainproduct.Service,
	priceHistoryWriter domainpricehistory.Writer,
	fetcher PriceFetcher,
	notifier PriceChangeNotifier,
) *PriceUpdater {
	return &PriceUpdater{
		productService:     productService,
		priceHistoryWriter: priceHistoryWriter,
		fetcher:            fetcher,
		notifier:           notifier,
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

	return u.updateProducts(ctx, market, products)
}

func (u *PriceUpdater) updateProducts(ctx context.Context, market enum.Market, products []domainproduct.Product) error {
	productMap := make(map[string]domainproduct.Product, len(products))
	externalIDs := make([]string, len(products))
	for i, p := range products {
		productMap[p.ExternalProductID] = p
		externalIDs[i] = p.ExternalProductID
	}

	for _, currency := range supportedCurrencies {
		prices, err := u.fetcher.FetchPrices(ctx, externalIDs, currency)
		if err != nil {
			log.Printf("failed to fetch prices for market=%s currency=%s: %v", market, currency, err)
			continue
		}

		updated := 0
		for _, pr := range prices {
			product, ok := productMap[pr.ExternalProductID]
			if !ok {
				continue
			}
			if pr.Price == "" {
				continue
			}

			if u.notifier != nil && u.notifySKUPriceChange(ctx, product, pr.Price, pr.Currency) {
				updated++
			}
		}

		log.Printf("batch price update: market=%s currency=%s total=%d fetched=%d updated=%d", market, currency, len(products), len(prices), updated)
	}

	return nil
}

func (u *PriceUpdater) notifySKUPriceChange(ctx context.Context, product domainproduct.Product, newPrice string, currency string) bool {
	if u.productService == nil {
		return false
	}

	skus, err := u.productService.FindSKUsByProductID(ctx, product.ID)
	if err != nil {
		log.Printf("failed to find skus for product %s: %v", product.ID, err)
		return false
	}

	if len(skus) != 1 {
		if len(skus) > 1 {
			log.Printf("skip sku price alert for product %s: sku-level diff source is required for multi-sku products", product.ID)
		}
		return false
	}

	sku := skus[0]
	if sku.Currency != "" && currency != "" && sku.Currency != currency {
		return false
	}

	oldPrice := sku.Price
	if oldPrice == "" || oldPrice == newPrice {
		return false
	}

	u.notifier.NotifyPriceChange(ctx, sku.ID, product.Title, oldPrice, newPrice, currency)
	return true
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
