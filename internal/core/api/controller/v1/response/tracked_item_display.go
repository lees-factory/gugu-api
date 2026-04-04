package response

import (
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type trackedItemDisplay struct {
	title        string
	mainImageURL string
	productURL   string
}

func resolveTrackedItemDisplay(product domainproduct.Product, variant *domainproduct.Variant) trackedItemDisplay {
	display := trackedItemDisplay{
		title:        product.Title,
		mainImageURL: product.MainImageURL,
		productURL:   product.ProductURL,
	}
	if variant == nil {
		return display
	}
	if variant.Title != "" {
		display.title = variant.Title
	}
	if variant.MainImageURL != "" {
		display.mainImageURL = variant.MainImageURL
	}
	if variant.ProductURL != "" {
		display.productURL = variant.ProductURL
	}
	return display
}
