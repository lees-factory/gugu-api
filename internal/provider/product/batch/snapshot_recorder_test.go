package batch

import (
	"context"
	"testing"
	"time"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	memorysnapshot "github.com/ljj/gugu-api/internal/storage/memory/pricesnapshot"
)

func TestRecordDailySnapshots_PublishesOnlyCompletedRuns(t *testing.T) {
	today := time.Date(2026, 4, 11, 0, 0, 0, 0, time.UTC)

	productService := newTestProductService(t,
		[]domainproduct.Product{
			{ID: "p-complete", Market: enum.MarketAliExpress, ExternalProductID: "100"},
			{ID: "p-partial", Market: enum.MarketAliExpress, ExternalProductID: "200"},
		},
		[]domainproduct.SKU{
			{ID: "sku-1", ProductID: "p-complete", Price: "100", OriginalPrice: "120", Currency: "USD"},
			{ID: "sku-2", ProductID: "p-complete", Price: "200", OriginalPrice: "220", Currency: "USD"},
			{ID: "sku-3", ProductID: "p-partial", Price: "", OriginalPrice: "320", Currency: "USD"},
			{ID: "sku-4", ProductID: "p-partial", Price: "300", OriginalPrice: "320", Currency: "USD"},
		},
	)

	skuRepo := memorysnapshot.NewSKUSnapshotRepository()
	recorder := NewSnapshotRecorder(productService, skuRepo)
	recorder.clock = func() time.Time { return today }

	if err := recorder.RecordDailySnapshots(context.Background()); err != nil {
		t.Fatalf("RecordDailySnapshots() error = %v", err)
	}

	complete1, _ := skuRepo.ListBySKUIDAndDateRange(context.Background(), "sku-1", "USD", today, today)
	complete2, _ := skuRepo.ListBySKUIDAndDateRange(context.Background(), "sku-2", "USD", today, today)
	partial1, _ := skuRepo.ListBySKUIDAndDateRange(context.Background(), "sku-3", "USD", today, today)
	partial2, _ := skuRepo.ListBySKUIDAndDateRange(context.Background(), "sku-4", "USD", today, today)

	if len(complete1) != 1 || len(complete2) != 1 {
		t.Fatalf("completed product snapshots were not fully published: sku-1=%d sku-2=%d", len(complete1), len(complete2))
	}
	if len(partial1) != 0 || len(partial2) != 0 {
		t.Fatalf("partial product snapshots must not be published: sku-3=%d sku-4=%d", len(partial1), len(partial2))
	}
}
