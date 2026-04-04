package request

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type GetTrackedItemPriceAlert struct {
	User          auth.RequestUser
	TrackedItemID string
	SKUID         string
}

func ParseGetTrackedItemPriceAlert(r *http.Request) GetTrackedItemPriceAlert {
	return GetTrackedItemPriceAlert{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
		SKUID:         strings.TrimSpace(r.URL.Query().Get("sku_id")),
	}
}
