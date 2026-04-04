package trackeditem

import "testing"

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
