package request

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type RegisterPriceAlert struct {
	User      auth.RequestUser
	ProductID string
	Channel   string
}

func ParseRegisterPriceAlert(r *http.Request) (RegisterPriceAlert, error) {
	var body struct {
		Channel string `json:"channel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return RegisterPriceAlert{}, err
	}
	return RegisterPriceAlert{
		User:      auth.RequestUserFrom(r.Context()),
		ProductID: strings.TrimSpace(chi.URLParam(r, "productID")),
		Channel:   body.Channel,
	}, nil
}

type UnregisterPriceAlert struct {
	User      auth.RequestUser
	ProductID string
}

func ParseUnregisterPriceAlert(r *http.Request) UnregisterPriceAlert {
	return UnregisterPriceAlert{
		User:      auth.RequestUserFrom(r.Context()),
		ProductID: strings.TrimSpace(chi.URLParam(r, "productID")),
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
