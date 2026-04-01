package batch

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	provideraliexpress "github.com/ljj/gugu-api/internal/provider/product/aliexpress"
)

var supportedCurrencies = []string{"KRW", "USD"}

type PriceChangeNotifier interface {
	NotifyPriceChange(ctx context.Context, productID string, productTitle string, oldPrice string, newPrice string, currency string)
}

type PriceFetcher interface {
	FetchPrices(ctx context.Context, externalProductIDs []string, currency string) ([]provideraliexpress.PriceResult, error)
}

type PriceUpdater struct {
	productService        *domainproduct.Service
	priceHistoryWriter    domainpricehistory.Writer
	productSnapshotWriter domainps.ProductSnapshotWriter
	fetcher               PriceFetcher
	notifier              PriceChangeNotifier
	clock                 func() time.Time
}

func NewPriceUpdater(
	productService *domainproduct.Service,
	priceHistoryWriter domainpricehistory.Writer,
	productSnapshotWriter domainps.ProductSnapshotWriter,
	fetcher PriceFetcher,
	notifier PriceChangeNotifier,
) *PriceUpdater {
	return &PriceUpdater{
		productService:        productService,
		priceHistoryWriter:    priceHistoryWriter,
		productSnapshotWriter: productSnapshotWriter,
		fetcher:               fetcher,
		notifier:              notifier,
		clock:                 func() time.Time { return time.Now() },
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

	now := u.clock()
	today := truncateToDate(now)

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

			if currency == product.Currency {
				if product.CurrentPrice != pr.Price {
					changeValue := calculateChange(product.CurrentPrice, pr.Price)

					if u.productService != nil {
						if err := u.productService.UpdatePrice(ctx, product.ID, pr.Price, pr.Currency); err != nil {
							log.Printf("failed to update price for product %s: %v", product.ID, err)
							continue
						}
					}

					if err := u.priceHistoryWriter.Create(ctx, domainpricehistory.PriceHistory{
						ProductID:   product.ID,
						Price:       pr.Price,
						Currency:    pr.Currency,
						RecordedAt:  now,
						ChangeValue: changeValue,
					}); err != nil {
						log.Printf("failed to record price history for product %s currency %s: %v", product.ID, currency, err)
					}

					if u.notifier != nil {
						u.notifier.NotifyPriceChange(ctx, product.ID, product.Title, product.CurrentPrice, pr.Price, pr.Currency)
					}

					updated++
				}
			} else {
				if err := u.priceHistoryWriter.Create(ctx, domainpricehistory.PriceHistory{
					ProductID:   product.ID,
					Price:       pr.Price,
					Currency:    currency,
					RecordedAt:  now,
					ChangeValue: "0",
				}); err != nil {
					log.Printf("failed to record price history for product %s currency %s: %v", product.ID, currency, err)
				}
				updated++
			}

			if err := u.productSnapshotWriter.Upsert(ctx, domainps.ProductPriceSnapshot{
				ProductID:    product.ID,
				SnapshotDate: today,
				Price:        pr.Price,
				Currency:     currency,
			}); err != nil {
				log.Printf("failed to upsert product snapshot for product %s currency %s: %v", product.ID, currency, err)
			}
		}

		log.Printf("batch price update: market=%s currency=%s total=%d fetched=%d updated=%d", market, currency, len(products), len(prices), updated)
	}

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
