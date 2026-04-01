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
)

func TestUpdateProducts_MultiCurrency(t *testing.T) {
	fixedTime := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)

	products := []domainproduct.Product{
		{ID: "p1", Market: enum.MarketAliExpress, ExternalProductID: "100", CurrentPrice: "7700", Currency: "KRW"},
	}

	historyWriter := &stubPriceHistoryWriter{}
	snapshotWriter := &stubProductSnapshotWriter{}

	updater := &PriceUpdater{
		productService:        nil,
		priceHistoryWriter:    historyWriter,
		productSnapshotWriter: snapshotWriter,
		fetcher: &stubFetcher{
			resultsByCurrency: map[string][]provideraliexpress.PriceResult{
				"KRW": {{ExternalProductID: "100", Price: "7800", Currency: "KRW"}},
				"USD": {{ExternalProductID: "100", Price: "5.99", Currency: "USD"}},
			},
		},
		notifier: &stubNotifier{},
		clock:    func() time.Time { return fixedTime },
	}

	err := updater.updateProducts(context.Background(), enum.MarketAliExpress, products)
	if err != nil {
		t.Fatalf("updateProducts() error = %v", err)
	}

	// KRW history (price changed 7700 -> 7800) + USD history
	if len(historyWriter.histories) != 2 {
		t.Fatalf("expected 2 history records, got %d", len(historyWriter.histories))
	}

	krwFound := false
	usdFound := false
	for _, h := range historyWriter.histories {
		if h.Currency == "KRW" && h.Price == "7800" {
			krwFound = true
			if h.ChangeValue != "100" {
				t.Errorf("KRW change_value = %q, want 100", h.ChangeValue)
			}
		}
		if h.Currency == "USD" && h.Price == "5.99" {
			usdFound = true
		}
	}
	if !krwFound {
		t.Error("KRW price history not found")
	}
	if !usdFound {
		t.Error("USD price history not found")
	}

	// Snapshots for both currencies
	if len(snapshotWriter.snapshots) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(snapshotWriter.snapshots))
	}
}

func TestUpdateProducts_NoPriceChange_SkipsRepresentativeHistory(t *testing.T) {
	fixedTime := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)

	products := []domainproduct.Product{
		{ID: "p1", Market: enum.MarketAliExpress, ExternalProductID: "100", CurrentPrice: "7700", Currency: "KRW"},
	}

	historyWriter := &stubPriceHistoryWriter{}
	snapshotWriter := &stubProductSnapshotWriter{}

	updater := &PriceUpdater{
		priceHistoryWriter:    historyWriter,
		productSnapshotWriter: snapshotWriter,
		fetcher: &stubFetcher{
			resultsByCurrency: map[string][]provideraliexpress.PriceResult{
				"KRW": {{ExternalProductID: "100", Price: "7700", Currency: "KRW"}}, // same price
				"USD": {{ExternalProductID: "100", Price: "5.99", Currency: "USD"}},
			},
		},
		notifier: &stubNotifier{},
		clock:    func() time.Time { return fixedTime },
	}

	err := updater.updateProducts(context.Background(), enum.MarketAliExpress, products)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	// KRW price unchanged -> no KRW history, only USD
	if len(historyWriter.histories) != 1 {
		t.Fatalf("expected 1 history record (USD only), got %d", len(historyWriter.histories))
	}
	if historyWriter.histories[0].Currency != "USD" {
		t.Errorf("expected USD history, got %s", historyWriter.histories[0].Currency)
	}

	// Snapshots still both (upsert always)
	if len(snapshotWriter.snapshots) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(snapshotWriter.snapshots))
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

type stubNotifier struct{}

func (n *stubNotifier) NotifyPriceChange(_ context.Context, _ string, _ string, _ string, _ string, _ string) {
}
