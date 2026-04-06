package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRegisterOpenAPIRouteSetsNoCacheHeaders(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	if err := os.Chdir("../../../"); err != nil {
		t.Fatalf("chdir to repo root: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()

	router := chi.NewRouter()
	registerOpenAPIRoute(router)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yml", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if got := rec.Header().Get("Cache-Control"); got != "no-cache, no-store, must-revalidate" {
		t.Fatalf("expected Cache-Control header to be set, got %q", got)
	}

	if got := rec.Header().Get("Pragma"); got != "no-cache" {
		t.Fatalf("expected Pragma header to be set, got %q", got)
	}

	if got := rec.Header().Get("Expires"); got != "0" {
		t.Fatalf("expected Expires header to be set, got %q", got)
	}
}
