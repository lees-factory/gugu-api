package trackeditem

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	memorypricealert "github.com/ljj/gugu-api/internal/storage/memory/pricealert"
)

func TestResolveSKUPriceHistoryCurrency_UsesOverride(t *testing.T) {
	got := resolveSKUPriceHistoryCurrency("usd", "KRW")
	if got != "USD" {
		t.Fatalf("resolveSKUPriceHistoryCurrency() = %q, want USD", got)
	}
}

func TestResolveSKUPriceHistoryCurrency_UsesTrackedItemCurrencyByDefault(t *testing.T) {
	got := resolveSKUPriceHistoryCurrency("", "krw")
	if got != "KRW" {
		t.Fatalf("resolveSKUPriceHistoryCurrency() = %q, want KRW", got)
	}
}

func TestResolveSKUPriceHistoryCurrency_FallsBackToKRW(t *testing.T) {
	got := resolveSKUPriceHistoryCurrency("", "")
	if got != "KRW" {
		t.Fatalf("resolveSKUPriceHistoryCurrency() = %q, want KRW", got)
	}
}

func TestResolvePriceAlertStateBySKUID_DefaultsOffWithoutService(t *testing.T) {
	controller := &Controller{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tracked-items/tracked-1/price-alert", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}

	got := controller.resolvePriceAlertStateBySKUID(req, "user-1", "sku-1")

	if got == nil {
		t.Fatal("resolvePriceAlertStateBySKUID() returned nil")
	}
	if *got != (response.PriceAlertState{Enabled: false}) {
		t.Fatalf("resolvePriceAlertStateBySKUID() = %+v, want disabled default state", *got)
	}
}

func TestResolvePriceAlertStateBySKUID_ReturnsStoredAlert(t *testing.T) {
	repo := memorypricealert.NewRepository()
	service := domainpricealert.NewService(
		domainpricealert.NewFinder(repo),
		repo,
		testAlertIDGenerator{},
		testAlertClock{},
	)
	if _, err := service.Register(context.Background(), "user-1", "sku-1", "EMAIL"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	controller := &Controller{priceAlertService: service}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tracked-items/tracked-1/price-alert?sku_id=sku-1", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}

	got := controller.resolvePriceAlertStateBySKUID(req, "user-1", "sku-1")

	if got == nil {
		t.Fatal("resolvePriceAlertStateBySKUID() returned nil")
	}
	if !got.Enabled {
		t.Fatalf("resolvePriceAlertStateBySKUID().Enabled = false, want true")
	}
	if got.Channel != "EMAIL" {
		t.Fatalf("resolvePriceAlertStateBySKUID().Channel = %q, want EMAIL", got.Channel)
	}
}

func TestResolveTrackedItemPriceAlertSKUID_UsesRequestedSKU(t *testing.T) {
	detail := &domaintrackeditem.TrackedItemDetail{
		TrackedItem: domaintrackeditem.TrackedItem{SKUID: "sku-selected"},
		SKUs: []domainproduct.SKU{
			{ID: "sku-selected"},
			{ID: "sku-requested"},
		},
	}

	got, err := resolveTrackedItemPriceAlertSKUID(detail, "sku-requested", true)
	if err != nil {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() error = %v", err)
	}
	if got != "sku-requested" {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() = %q, want sku-requested", got)
	}
}

func TestResolveTrackedItemPriceAlertSKUID_UsesSelectedSKUByDefault(t *testing.T) {
	detail := &domaintrackeditem.TrackedItemDetail{
		TrackedItem: domaintrackeditem.TrackedItem{SKUID: "sku-selected"},
		SKUs:        []domainproduct.SKU{{ID: "sku-selected"}},
	}

	got, err := resolveTrackedItemPriceAlertSKUID(detail, "", true)
	if err != nil {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() error = %v", err)
	}
	if got != "sku-selected" {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() = %q, want sku-selected", got)
	}
}

func TestResolveTrackedItemPriceAlertSKUID_AllowsMissingSKUForRead(t *testing.T) {
	detail := &domaintrackeditem.TrackedItemDetail{}

	got, err := resolveTrackedItemPriceAlertSKUID(detail, "", false)
	if err != nil {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() error = %v", err)
	}
	if got != "" {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() = %q, want empty", got)
	}
}

type testAlertIDGenerator struct{}

func (testAlertIDGenerator) New() (string, error) {
	return "alert-1", nil
}

type testAlertClock struct{}

func (testAlertClock) Now() time.Time {
	return time.Unix(0, 0)
}
