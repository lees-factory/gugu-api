package request

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type GetSKUPriceHistories struct {
	User          auth.RequestUser
	TrackedItemID string
	SKUID         string
	Currency      string
}

func ParseGetSKUPriceHistories(r *http.Request) GetSKUPriceHistories {
	return GetSKUPriceHistories{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
		SKUID:         strings.TrimSpace(r.URL.Query().Get("sku_id")),
		Currency:      strings.TrimSpace(r.URL.Query().Get("currency")),
	}
}
