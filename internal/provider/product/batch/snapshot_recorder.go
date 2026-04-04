package batch

import (
	"context"
	"fmt"
	"log"
	"time"

	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
)

type SnapshotRecorder struct {
	productService        *domainproduct.Service
	productSnapshotWriter domainps.ProductSnapshotWriter
	skuSnapshotWriter     domainps.SKUSnapshotWriter
	clock                 func() time.Time
}

func NewSnapshotRecorder(
	productService *domainproduct.Service,
	productSnapshotWriter domainps.ProductSnapshotWriter,
	skuSnapshotWriter domainps.SKUSnapshotWriter,
) *SnapshotRecorder {
	return &SnapshotRecorder{
		productService:        productService,
		productSnapshotWriter: productSnapshotWriter,
		skuSnapshotWriter:     skuSnapshotWriter,
		clock:                 func() time.Time { return time.Now() },
	}
}

func (r *SnapshotRecorder) RecordDailySnapshots(ctx context.Context) error {
	today := truncateToDate(r.clock())

	products, err := r.productService.ListByMarket(ctx, enum.MarketAliExpress)
	if err != nil {
		return fmt.Errorf("list products: %w", err)
	}

	productRecorded := 0
	skuRecorded := 0

	for _, p := range products {
		skus, err := r.productService.FindSKUsByProductID(ctx, p.ID)
		if err != nil {
			log.Printf("failed to find skus for product %s: %v", p.ID, err)
			continue
		}

		for _, sku := range skus {
			if sku.Price == "" {
				continue
			}
			if err := r.skuSnapshotWriter.Upsert(ctx, domainps.SKUPriceSnapshot{
				SKUID:         sku.ID,
				SnapshotDate:  today,
				Price:         sku.Price,
				OriginalPrice: sku.OriginalPrice,
				Currency:      sku.Currency,
			}); err != nil {
				log.Printf("failed to record sku snapshot for %s: %v", sku.ID, err)
				continue
			}
			skuRecorded++
		}
		if len(skus) > 0 {
			productRecorded++
		}
	}

	log.Printf("daily snapshot recorded: products=%d skus=%d", productRecorded, skuRecorded)
	return nil
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
