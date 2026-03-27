package request

import (
	"encoding/json"
	"net/http"

	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type AddTrackedItem struct {
	User              auth.RequestUser
	ProviderCommerce  string
	ExternalProductID string
	OriginalURL       string
}

func ParseAddTrackedItem(r *http.Request) (AddTrackedItem, error) {
	var body struct {
		ProviderCommerce  string `json:"provider_commerce"`
		ExternalProductID string `json:"external_product_id"`
		OriginalURL       string `json:"original_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return AddTrackedItem{}, err
	}
	return AddTrackedItem{
		User:              auth.RequestUserFrom(r.Context()),
		ProviderCommerce:  body.ProviderCommerce,
		ExternalProductID: body.ExternalProductID,
		OriginalURL:       body.OriginalURL,
	}, nil
}
