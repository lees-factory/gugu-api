package request

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseAddTrackedItems_ParsesLanguage(t *testing.T) {
	req := httptest.NewRequest("POST", "/v1/tracked-items/", strings.NewReader(`{
		"items": [
			{
				"provider_commerce": "ALIEXPRESS",
				"origin_product_id": "3256809920794713",
				"external_product_id": "1001",
				"original_url": "https://example.com/item/1001",
				"currency": "KRW",
				"language": "KO"
			}
		]
	}`))

	parsed, err := ParseAddTrackedItems(req)
	if err != nil {
		t.Fatalf("ParseAddTrackedItems() error = %v", err)
	}
	if len(parsed.Items) != 1 {
		t.Fatalf("items count = %d, want 1", len(parsed.Items))
	}
	if parsed.Items[0].Language != "KO" {
		t.Fatalf("language = %q, want KO", parsed.Items[0].Language)
	}
	if parsed.Items[0].OriginProductID != "3256809920794713" {
		t.Fatalf("origin_product_id = %q, want 3256809920794713", parsed.Items[0].OriginProductID)
	}
}
