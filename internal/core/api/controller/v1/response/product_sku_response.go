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

func NewProductSKUs(skus []domainproduct.SKU) []ProductSKU {
	result := make([]ProductSKU, len(skus))
	for i, sku := range skus {
		item := newProductSKU(sku)
		item.CurrentPrice = sku.Price
		result[i] = item
	}
	return result
}

func NewProductSKUsWithCurrentPrice(skus []domainproduct.SKU, currentPriceBySKUID map[string]string) []ProductSKU {
	result := make([]ProductSKU, len(skus))
	for i, sku := range skus {
		item := newProductSKU(sku)
		item.CurrentPrice = strings.TrimSpace(currentPriceBySKUID[strings.TrimSpace(sku.ID)])
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
		Price:         sku.Price,
		OriginalPrice: sku.OriginalPrice,
		Currency:      sku.Currency,
		ImageURL:      sku.ImageURL,
		SKUProperties: sku.SKUProperties,
	}
}
