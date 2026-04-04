package request

import (
	"encoding/json"
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

type RegisterTrackedItemPriceAlert struct {
	User          auth.RequestUser
	TrackedItemID string
	SKUID         string
	Channel       string
}

func ParseRegisterTrackedItemPriceAlert(r *http.Request) (RegisterTrackedItemPriceAlert, error) {
	var body struct {
		SKUID   string `json:"sku_id"`
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return RegisterTrackedItemPriceAlert{}, err
	}

	return RegisterTrackedItemPriceAlert{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
		SKUID:         strings.TrimSpace(body.SKUID),
		Channel:       strings.TrimSpace(body.Channel),
	}, nil
}

type UnregisterTrackedItemPriceAlert struct {
	User          auth.RequestUser
	TrackedItemID string
	SKUID         string
}

func ParseUnregisterTrackedItemPriceAlert(r *http.Request) UnregisterTrackedItemPriceAlert {
	return UnregisterTrackedItemPriceAlert{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
		SKUID:         strings.TrimSpace(r.URL.Query().Get("sku_id")),
	}
}
