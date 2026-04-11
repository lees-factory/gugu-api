package response

import (
	"strings"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type ProductSKU struct {
	ID            string `json:"id"`
	ExternalSKUID string `json:"external_sku_id"`
	OriginSKUID   string `json:"origin_sku_id,omitempty"`
	SKUName       string `json:"sku_name"`
	Color         string `json:"color"`
	Size          string `json:"size"`
	Price         string `json:"price"`
	CurrentPrice  string `json:"current_price"`
	OriginalPrice string `json:"original_price"`
	Currency      string `json:"currency"`
	ImageURL      string `json:"image_url"`
	SKUProperties string `json:"sku_properties,omitempty"`
}

type SKUCurrentSnapshot struct {
	Price         string
	OriginalPrice string
	Currency      string
}

func NewProductSKUsWithCurrentPrice(skus []domainproduct.SKU, currentBySKUID map[string]SKUCurrentSnapshot) []ProductSKU {
	result := make([]ProductSKU, len(skus))
	for i, sku := range skus {
		item := newProductSKU(sku)
		current := currentBySKUID[strings.TrimSpace(sku.ID)]
		item.CurrentPrice = strings.TrimSpace(current.Price)
		item.Price = strings.TrimSpace(current.Price)
		item.OriginalPrice = strings.TrimSpace(current.OriginalPrice)
		item.Currency = strings.TrimSpace(current.Currency)
		result[i] = item
	}
	return result
}

func newProductSKU(sku domainproduct.SKU) ProductSKU {
	return ProductSKU{
		ID:            sku.ID,
		ExternalSKUID: sku.ExternalSKUID,
		OriginSKUID:   sku.OriginSKUID,
		SKUName:       sku.SKUName,
		Color:         sku.Color,
		Size:          sku.Size,
		Price:         "",
		OriginalPrice: "",
		Currency:      "",
		ImageURL:      sku.ImageURL,
		SKUProperties: sku.SKUProperties,
	}
}
