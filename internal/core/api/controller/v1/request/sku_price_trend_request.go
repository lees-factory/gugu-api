package request

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type GetSKUPriceTrend struct {
	User          auth.RequestUser
	TrackedItemID string
	SKUID         string
	Currency      string
	From          string
	To            string
}

func ParseGetSKUPriceTrend(r *http.Request) GetSKUPriceTrend {
	return GetSKUPriceTrend{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
		SKUID:         strings.TrimSpace(r.URL.Query().Get("sku_id")),
		Currency:      strings.TrimSpace(r.URL.Query().Get("currency")),
		From:          strings.TrimSpace(r.URL.Query().Get("from")),
		To:            strings.TrimSpace(r.URL.Query().Get("to")),
	}
}
