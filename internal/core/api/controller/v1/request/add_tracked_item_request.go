package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ljj/gugu-api/internal/core/support/auth"
)

const maxAddTrackedItems = 5

type AddTrackedItemEntry struct {
	ProviderCommerce  string `json:"provider_commerce"`
	ExternalProductID string `json:"external_product_id"`
	OriginalURL       string `json:"original_url"`
	Currency          string `json:"currency"`
	Language          string `json:"language"`
}

type AddTrackedItems struct {
	User  auth.RequestUser
	Items []AddTrackedItemEntry
}

func ParseAddTrackedItems(r *http.Request) (AddTrackedItems, error) {
	var body struct {
		Items []AddTrackedItemEntry `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return AddTrackedItems{}, err
	}

	if len(body.Items) == 0 {
		return AddTrackedItems{}, fmt.Errorf("items must not be empty")
	}
	if len(body.Items) > maxAddTrackedItems {
		return AddTrackedItems{}, fmt.Errorf("items must not exceed %d", maxAddTrackedItems)
	}

	return AddTrackedItems{
		User:  auth.RequestUserFrom(r.Context()),
		Items: body.Items,
	}, nil
}
