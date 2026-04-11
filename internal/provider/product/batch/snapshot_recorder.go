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
	productService    *domainproduct.Service
	skuSnapshotWriter domainps.SKUSnapshotWriter
	clock             func() time.Time
}

func NewSnapshotRecorder(
	productService *domainproduct.Service,
	skuSnapshotWriter domainps.SKUSnapshotWriter,
) *SnapshotRecorder {
	return &SnapshotRecorder{
		productService:    productService,
		skuSnapshotWriter: skuSnapshotWriter,
		clock:             func() time.Time { return time.Now() },
	}
}

func (r *SnapshotRecorder) RecordDailySnapshots(ctx context.Context) error {
	today := truncateToDate(r.clock())

	products, err := r.productService.ListByMarket(ctx, enum.MarketAliExpress)
	if err != nil {
		return fmt.Errorf("list products: %w", err)
	}

	totalRuns := 0
	completedRuns := 0
	partialRuns := 0
	failedRuns := 0
	productRecorded := 0
	skuRecorded := 0

	for _, p := range products {
		totalRuns++
		skus, err := r.productService.FindSKUsByProductID(ctx, p.ID)
		if err != nil {
			log.Printf("failed to find skus for product %s: %v", p.ID, err)
			failedRuns++
			continue
		}

		expectedSKUs := len(skus)
		if expectedSKUs == 0 {
			partialRuns++
			log.Printf("snapshot run partial: product=%s expected_sku_count=0 collected_sku_count=0", p.ID)
			continue
		}

		candidates := make([]domainps.SKUPriceSnapshot, 0, len(skus))
		for _, sku := range skus {
			if sku.Price == "" || sku.Currency == "" {
				continue
			}
			candidates = append(candidates, domainps.SKUPriceSnapshot{
				SKUID:         sku.ID,
				SnapshotDate:  today,
				Price:         sku.Price,
				OriginalPrice: sku.OriginalPrice,
				Currency:      sku.Currency,
			})
		}

		collectedSKUs := len(candidates)
		if collectedSKUs != expectedSKUs {
			partialRuns++
			log.Printf("snapshot run partial: product=%s expected_sku_count=%d collected_sku_count=%d", p.ID, expectedSKUs, collectedSKUs)
			continue
		}

		publishFailed := false
		for _, snapshot := range candidates {
			if err := r.skuSnapshotWriter.Upsert(ctx, snapshot); err != nil {
				log.Printf("failed to record sku snapshot for %s: %v", snapshot.SKUID, err)
				failedRuns++
				publishFailed = true
				break
			}
			skuRecorded++
		}
		if publishFailed {
			continue
		}

		completedRuns++
		productRecorded++
	}

	completionRate := 0.0
	if totalRuns > 0 {
		completionRate = float64(completedRuns) / float64(totalRuns)
	}

	if partialRuns > 0 || failedRuns > 0 {
		log.Printf("snapshot run alert: completion_rate=%.4f total_runs=%d completed_runs=%d partial_runs=%d failed_runs=%d",
			completionRate, totalRuns, completedRuns, partialRuns, failedRuns)
	} else {
		log.Printf("snapshot run summary: completion_rate=%.4f total_runs=%d completed_runs=%d partial_runs=%d failed_runs=%d",
			completionRate, totalRuns, completedRuns, partialRuns, failedRuns)
	}

	log.Printf("daily snapshot recorded: products=%d skus=%d", productRecorded, skuRecorded)
	return nil
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
