package request

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type SelectSKU struct {
	User          auth.RequestUser
	TrackedItemID string
	SKUID         string
}

func ParseSelectSKU(r *http.Request) (SelectSKU, error) {
	var body struct {
		SKUID string `json:"sku_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return SelectSKU{}, err
	}
	return SelectSKU{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
		SKUID:         body.SKUID,
	}, nil
}
