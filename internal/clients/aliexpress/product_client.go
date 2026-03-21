package aliexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (c *HTTPClient) GetProductSnapshot(ctx context.Context, input ProductLookupInput) (*ProductSnapshot, error) {
	result, err := c.GetAffiliateProductDetail(ctx, ProductDetailInput{
		ProductIDs:     []string{input.ProductID},
		TargetCurrency: input.TargetCurrency,
		TargetLanguage: input.TargetLanguage,
		Country:        input.Country,
		TrackingID:     input.TrackingID,
		Fields:         input.Fields,
	})
	if err != nil {
		return nil, err
	}
	if len(result.Products) == 0 {
		return nil, fmt.Errorf("aliexpress product detail returned no products")
	}

	product := result.Products[0]
	price := firstNonEmpty(product.TargetSalePrice, product.SalePrice, product.TargetAppSalePrice, product.AppSalePrice)
	currency := firstNonEmpty(product.TargetSalePriceCurrency, product.SalePriceCurrency, product.TargetAppSalePriceCurrency, product.AppSalePriceCurrency)

	return &ProductSnapshot{
		ProductID:      strconv.FormatInt(product.ProductID, 10),
		Title:          product.ProductTitle,
		Price:          price,
		Currency:       currency,
		MainImageURL:   product.ProductMainImageURL,
		ProductURL:     product.ProductDetailURL,
		PromotionLink:  product.PromotionLink,
		OriginalPrice:  firstNonEmpty(product.TargetOriginalPrice, product.OriginalPrice),
		StoreName:      product.ShopName,
		TrackingIDUsed: strings.TrimSpace(input.TrackingID),
	}, nil
}

func (c *HTTPClient) GetAffiliateProductDetail(ctx context.Context, input ProductDetailInput) (*ProductDetailResult, error) {
	productIDs := normalizeCSV(input.ProductIDs)
	if len(productIDs) == 0 {
		return nil, fmt.Errorf("product ids are required")
	}

	fields := input.Fields
	if len(fields) == 0 {
		fields = []string{
			"commission_rate",
			"sale_price",
			"target_sale_price",
			"target_sale_price_currency",
			"product_main_image_url",
			"product_title",
			"promotion_link",
			"product_detail_url",
			"target_original_price",
			"shop_name",
			"product_small_image_urls",
			"sku_id",
		}
	}

	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName: "aliexpress.affiliate.productdetail.get",
		form: map[string]string{
			"fields":          strings.Join(fields, ","),
			"product_ids":     strings.Join(productIDs, ","),
			"target_currency": defaultString(input.TargetCurrency, "USD"),
			"target_language": defaultString(input.TargetLanguage, "EN"),
			"country":         defaultString(input.Country, "US"),
			"tracking_id":     strings.TrimSpace(input.TrackingID),
		},
		topProtocol: true,
	})
	if err != nil {
		return nil, err
	}

	var payload productDetailEnvelope
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode product detail response: %w", err)
	}
	if payload.RespResult.RespCode != 20010000 {
		return &ProductDetailResult{}, nil
	}

	return &ProductDetailResult{
		CurrentRecordCount: payload.RespResult.Result.CurrentRecordCount,
		Products:           payload.RespResult.Result.Products.Product,
	}, nil
}

func (c *HTTPClient) GetAffiliateProductSKUDetail(ctx context.Context, input ProductSKUDetailInput) (*ProductSKUDetailResult, error) {
	productID := strings.TrimSpace(input.ProductID)
	if productID == "" {
		return nil, fmt.Errorf("product id is required")
	}
	shipToCountry := strings.TrimSpace(input.ShipToCountry)
	if shipToCountry == "" {
		return nil, fmt.Errorf("ship to country is required")
	}

	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName: "aliexpress.affiliate.product.sku.detail.get",
		form: map[string]string{
			"product_id":        productID,
			"ship_to_country":   shipToCountry,
			"target_currency":   defaultString(input.TargetCurrency, "USD"),
			"target_language":   defaultString(input.TargetLanguage, "EN"),
			"need_deliver_info": defaultString(input.NeedDeliverInfo, "No"),
			"sku_ids":           strings.Join(normalizeCSV(input.SKUIDs), ","),
		},
		topProtocol: true,
	})
	if err != nil {
		return nil, err
	}

	var payload productSKUDetailEnvelope
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode product sku detail response: %w", err)
	}

	return &ProductSKUDetailResult{
		Code:     payload.Result.Result.Code,
		Success:  payload.Result.Result.Success,
		ItemInfo: payload.Result.Result.ItemInfo,
		SKUInfos: payload.Result.Result.SKUInfos.TrafficSKUInfoList,
	}, nil
}
