package discover

import (
	"net/http/httptest"
	"testing"
)

func TestParseHotProductLanguage_UsesLanguageFirst(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/discover/hot-products?language=en&target_language=ko", nil)

	if got := parseHotProductLanguage(req); got != "EN" {
		t.Fatalf("parseHotProductLanguage() = %q, want EN", got)
	}
}

func TestParseHotProductLanguage_FallsBackToTargetLanguage(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/discover/hot-products?target_language=ko", nil)

	if got := parseHotProductLanguage(req); got != "KO" {
		t.Fatalf("parseHotProductLanguage() = %q, want KO", got)
	}
}

func TestParseHotProductLanguage_DefaultsToKO(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/discover/hot-products", nil)

	if got := parseHotProductLanguage(req); got != "KO" {
		t.Fatalf("parseHotProductLanguage() = %q, want KO", got)
	}
}
