package request

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseUpdateTrackedItemLanguage(t *testing.T) {
	req := httptest.NewRequest("PATCH", "/v1/tracked-items/tracked-1/language", strings.NewReader(`{
		"language": "en"
	}`))

	parsed, err := ParseUpdateTrackedItemLanguage(req)
	if err != nil {
		t.Fatalf("ParseUpdateTrackedItemLanguage() error = %v", err)
	}
	if parsed.Language != "EN" {
		t.Fatalf("language = %q, want EN", parsed.Language)
	}
}
