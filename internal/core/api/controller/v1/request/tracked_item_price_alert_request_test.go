package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestParseGetTrackedItemPriceAlert_PrefersPathSKUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/skus/sku-path/price-alert?sku_id=sku-query", nil)
	req = withRouteParams(req, map[string]string{
		"skuID": "sku-path",
	})

	parsed := ParseGetTrackedItemPriceAlert(req)
	if parsed.SKUID != "sku-path" {
		t.Fatalf("sku_id = %q, want sku-path", parsed.SKUID)
	}
}

func TestParseRegisterTrackedItemPriceAlert_AllowsEmptyBodyWithPathSKUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/skus/sku-path/price-alert", strings.NewReader(""))
	req = withRouteParams(req, map[string]string{
		"skuID": "sku-path",
	})

	parsed, err := ParseRegisterTrackedItemPriceAlert(req)
	if err != nil {
		t.Fatalf("ParseRegisterTrackedItemPriceAlert() error = %v", err)
	}
	if parsed.SKUID != "sku-path" {
		t.Fatalf("sku_id = %q, want sku-path", parsed.SKUID)
	}
	if parsed.Channel != "" {
		t.Fatalf("channel = %q, want empty", parsed.Channel)
	}
}

func TestParseUnregisterTrackedItemPriceAlert_UsesQueryFallback(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/v1/skus/price-alert?sku_id=sku-query", nil)

	parsed := ParseUnregisterTrackedItemPriceAlert(req)
	if parsed.SKUID != "sku-query" {
		t.Fatalf("sku_id = %q, want sku-query", parsed.SKUID)
	}
}

func withRouteParams(req *http.Request, params map[string]string) *http.Request {
	routeCtx := chi.NewRouteContext()
	for k, v := range params {
		routeCtx.URLParams.Add(k, v)
	}
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	return req.WithContext(ctx)
}
