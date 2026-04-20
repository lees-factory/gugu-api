package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type GetTrackedItemPriceAlert struct {
	User  auth.RequestUser
	SKUID string
}

func ParseGetTrackedItemPriceAlert(r *http.Request) GetTrackedItemPriceAlert {
	return GetTrackedItemPriceAlert{
		User:  auth.RequestUserFrom(r.Context()),
		SKUID: resolveSKUID(chi.URLParam(r, "skuID"), r.URL.Query().Get("sku_id")),
	}
}

type RegisterTrackedItemPriceAlert struct {
	User    auth.RequestUser
	SKUID   string
	Channel string
}

func ParseRegisterTrackedItemPriceAlert(r *http.Request) (RegisterTrackedItemPriceAlert, error) {
	var body struct {
		SKUID   string `json:"sku_id"`
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil && !errors.Is(err, io.EOF) {
		return RegisterTrackedItemPriceAlert{}, err
	}

	pathSKUID := strings.TrimSpace(chi.URLParam(r, "skuID"))
	return RegisterTrackedItemPriceAlert{
		User:    auth.RequestUserFrom(r.Context()),
		SKUID:   resolveSKUID(pathSKUID, body.SKUID),
		Channel: strings.TrimSpace(body.Channel),
	}, nil
}

type UnregisterTrackedItemPriceAlert struct {
	User  auth.RequestUser
	SKUID string
}

func ParseUnregisterTrackedItemPriceAlert(r *http.Request) UnregisterTrackedItemPriceAlert {
	return UnregisterTrackedItemPriceAlert{
		User:  auth.RequestUserFrom(r.Context()),
		SKUID: resolveSKUID(chi.URLParam(r, "skuID"), r.URL.Query().Get("sku_id")),
	}
}

func resolveSKUID(primary string, fallback string) string {
	skuID := strings.TrimSpace(primary)
	if skuID != "" {
		return skuID
	}
	return strings.TrimSpace(fallback)
}
