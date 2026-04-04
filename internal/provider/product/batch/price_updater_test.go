package batch

import (
	"context"
	"testing"
	"time"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	provideraliexpress "github.com/ljj/gugu-api/internal/provider/product/aliexpress"
	memoryproduct "github.com/ljj/gugu-api/internal/storage/memory/product"
)

func TestUpdateProducts_NoProductService_NoNotification(t *testing.T) {
	fixedTime := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)

	products := []domainproduct.Product{
		{ID: "p1", Market: enum.MarketAliExpress, ExternalProductID: "100"},
	}

	notifier := &stubNotifier{}
	updater := &PriceUpdater{
		fetcher: &stubFetcher{
			resultsByCurrency: map[string][]provideraliexpress.PriceResult{
				"KRW": {{ExternalProductID: "100", Price: "7800", Currency: "KRW"}},
			},
		},
		notifier: notifier,
		clock:    func() time.Time { return fixedTime },
	}

	err := updater.updateProducts(context.Background(), enum.MarketAliExpress, products)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	if len(notifier.calls) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(notifier.calls))
	}
}

func TestCalculateChange(t *testing.T) {
	tests := []struct {
		old, new, want string
	}{
		{"7700", "7800", "100"},
		{"7800", "7700", "-100"},
		{"5.99", "5.99", "0"},
		{"", "100", "0"},
		{"100", "", "0"},
	}
	for _, tt := range tests {
		got := calculateChange(tt.old, tt.new)
		if got != tt.want {
			t.Errorf("calculateChange(%q, %q) = %q, want %q", tt.old, tt.new, got, tt.want)
		}
	}
}

func TestUpdateProducts_NotifiesSingleSKUOnly(t *testing.T) {
	fixedTime := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	productService := newTestProductService(t, []domainproduct.Product{
		{ID: "p1", Market: enum.MarketAliExpress, ExternalProductID: "100", Title: "Single SKU Product"},
	}, []domainproduct.SKU{
		{ID: "sku-1", ProductID: "p1", Price: "7700", Currency: "KRW"},
	})

	notifier := &stubNotifier{}
	updater := &PriceUpdater{
		productService:        productService,
		priceHistoryWriter:    &stubPriceHistoryWriter{},
		productSnapshotWriter: &stubProductSnapshotWriter{},
		fetcher: &stubFetcher{
			resultsByCurrency: map[string][]provideraliexpress.PriceResult{
				"KRW": {{ExternalProductID: "100", Price: "7800", Currency: "KRW"}},
			},
		},
		notifier: notifier,
		clock:    func() time.Time { return fixedTime },
	}

	if err := updater.updateProducts(context.Background(), enum.MarketAliExpress, []domainproduct.Product{
		{ID: "p1", Market: enum.MarketAliExpress, ExternalProductID: "100", Title: "Single SKU Product"},
	}); err != nil {
		t.Fatalf("updateProducts() error = %v", err)
	}

	if len(notifier.calls) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifier.calls))
	}
	if notifier.calls[0].skuID != "sku-1" {
		t.Fatalf("skuID = %q, want sku-1", notifier.calls[0].skuID)
	}
}

func TestUpdateProducts_SkipsNotificationForMultiSKUProduct(t *testing.T) {
	fixedTime := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	productService := newTestProductService(t, []domainproduct.Product{
		{ID: "p1", Market: enum.MarketAliExpress, ExternalProductID: "100", Title: "Multi SKU Product"},
	}, []domainproduct.SKU{
		{ID: "sku-1", ProductID: "p1", Price: "7700", Currency: "KRW"},
		{ID: "sku-2", ProductID: "p1", Price: "7600", Currency: "KRW"},
	})

	notifier := &stubNotifier{}
	updater := &PriceUpdater{
		productService:        productService,
		priceHistoryWriter:    &stubPriceHistoryWriter{},
		productSnapshotWriter: &stubProductSnapshotWriter{},
		fetcher: &stubFetcher{
			resultsByCurrency: map[string][]provideraliexpress.PriceResult{
				"KRW": {{ExternalProductID: "100", Price: "7800", Currency: "KRW"}},
			},
		},
		notifier: notifier,
		clock:    func() time.Time { return fixedTime },
	}

	if err := updater.updateProducts(context.Background(), enum.MarketAliExpress, []domainproduct.Product{
		{ID: "p1", Market: enum.MarketAliExpress, ExternalProductID: "100", Title: "Multi SKU Product"},
	}); err != nil {
		t.Fatalf("updateProducts() error = %v", err)
	}

	if len(notifier.calls) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(notifier.calls))
	}
}

// --- stubs ---

type stubFetcher struct {
	resultsByCurrency map[string][]provideraliexpress.PriceResult
}

func (f *stubFetcher) FetchPrices(_ context.Context, _ []string, currency string) ([]provideraliexpress.PriceResult, error) {
	return f.resultsByCurrency[currency], nil
}

type stubPriceHistoryWriter struct {
	histories []domainpricehistory.PriceHistory
}

func (w *stubPriceHistoryWriter) Create(_ context.Context, h domainpricehistory.PriceHistory) error {
	w.histories = append(w.histories, h)
	return nil
}

type stubProductSnapshotWriter struct {
	snapshots []domainps.ProductPriceSnapshot
}

func (w *stubProductSnapshotWriter) Upsert(_ context.Context, s domainps.ProductPriceSnapshot) error {
	w.snapshots = append(w.snapshots, s)
	return nil
}

type notifierCall struct {
	skuID string
}

type stubNotifier struct {
	calls []notifierCall
}

func (n *stubNotifier) NotifyPriceChange(_ context.Context, skuID string, _ string, _ string, _ string, _ string) {
	n.calls = append(n.calls, notifierCall{skuID: skuID})
}

type stubIDGenerator struct{}

func (g stubIDGenerator) New() (string, error) {
	return "generated-id", nil
}

type stubClock struct {
	now time.Time
}

func (c stubClock) Now() time.Time {
	return c.now
}

func newTestProductService(t *testing.T, products []domainproduct.Product, skus []domainproduct.SKU) *domainproduct.Service {
	t.Helper()

	productRepo := memoryproduct.NewRepository()
	skuRepo := memoryproduct.NewSKURepository()

	for _, p := range products {
		if err := productRepo.Create(context.Background(), p); err != nil {
			t.Fatalf("create product: %v", err)
		}
	}

	for _, sku := range skus {
		if err := skuRepo.Create(context.Background(), sku); err != nil {
			t.Fatalf("create sku: %v", err)
		}
	}

	return domainproduct.NewService(
		domainproduct.NewFinder(productRepo),
		domainproduct.NewWriter(productRepo),
		memoryproduct.NewVariantRepository(),
		skuRepo,
		stubIDGenerator{},
		stubClock{now: time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)},
		nil,
		nil,
	)
}
