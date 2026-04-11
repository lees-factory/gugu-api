package response

import (
	"strings"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type SKUHierarchy struct {
	ColorGroups []SKUColorGroup `json:"color_groups"`
}

type SKUColorGroup struct {
	Color    string         `json:"color"`
	ImageURL string         `json:"image_url"`
	Sizes    []SKUSizeEntry `json:"sizes"`
}

type SKUSizeEntry struct {
	Size          string `json:"size"`
	SKUID         string `json:"sku_id"`
	Price         string `json:"price"`
	OriginalPrice string `json:"original_price"`
	Currency      string `json:"currency"`
}

func BuildSKUHierarchy(skus []domainproduct.SKU, currentBySKUID map[string]SKUCurrentSnapshot) *SKUHierarchy {
	if len(skus) == 0 {
		return nil
	}

	colorOrder := make([]string, 0)
	colorMap := make(map[string]*SKUColorGroup)

	for _, sku := range skus {
		color := sku.Color
		if color == "" {
			color = "default"
		}

		group, exists := colorMap[color]
		if !exists {
			group = &SKUColorGroup{
				Color:    color,
				ImageURL: sku.ImageURL,
			}
			colorMap[color] = group
			colorOrder = append(colorOrder, color)
		}

		size := sku.Size
		if size == "" {
			size = "default"
		}

		group.Sizes = append(group.Sizes, SKUSizeEntry{
			Size:          size,
			SKUID:         sku.ID,
			Price:         strings.TrimSpace(currentBySKUID[strings.TrimSpace(sku.ID)].Price),
			OriginalPrice: strings.TrimSpace(currentBySKUID[strings.TrimSpace(sku.ID)].OriginalPrice),
			Currency:      strings.TrimSpace(currentBySKUID[strings.TrimSpace(sku.ID)].Currency),
		})
	}

	groups := make([]SKUColorGroup, 0, len(colorOrder))
	for _, color := range colorOrder {
		groups = append(groups, *colorMap[color])
	}

	return &SKUHierarchy{ColorGroups: groups}
}
