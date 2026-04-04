package request

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type RegisterPriceAlert struct {
	User    auth.RequestUser
	SKUID   string
	Channel string
}

func ParseRegisterPriceAlert(r *http.Request) (RegisterPriceAlert, error) {
	var body struct {
		Channel string `json:"channel"`
		SKUID   string `json:"sku_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return RegisterPriceAlert{}, err
	}
	return RegisterPriceAlert{
		User:    auth.RequestUserFrom(r.Context()),
		SKUID:   strings.TrimSpace(body.SKUID),
		Channel: body.Channel,
	}, nil
}

type UnregisterPriceAlert struct {
	User  auth.RequestUser
	SKUID string
}

func ParseUnregisterPriceAlert(r *http.Request) UnregisterPriceAlert {
	return UnregisterPriceAlert{
		User:  auth.RequestUserFrom(r.Context()),
		SKUID: strings.TrimSpace(chi.URLParam(r, "skuID")),
	}
}

type ListMyAlerts struct {
	User auth.RequestUser
}

func ParseListMyAlerts(r *http.Request) ListMyAlerts {
	return ListMyAlerts{
		User: auth.RequestUserFrom(r.Context()),
	}
}
