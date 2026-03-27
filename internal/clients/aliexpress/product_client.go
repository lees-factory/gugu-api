package aliexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (c *HTTPClient) GetAffiliateCategories(ctx context.Context) (*CategoryResult, error) {
	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName:     "aliexpress.affiliate.category.get",
		form:        map[string]string{},
		topProtocol: true,
	})
	if err != nil {
		return nil, err
	}

	var payload categoryEnvelope
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode category response: %w", err)
	}
	return &CategoryResult{
		TotalResultCount: payload.RespResult.Result.TotalResultCount,
		Categories:       payload.RespResult.Result.Categories.Category,
	}, nil
}

func (c *HTTPClient) GetAffiliateProducts(ctx context.Context, input ProductQueryInput) (*ProductQueryResult, error) {
	form := map[string]string{
		"target_currency": defaultString(input.TargetCurrency, "USD"),
		"target_language": defaultString(input.TargetLanguage, "EN"),
		"access_token":    strings.TrimSpace(input.AccessToken),
	}
	if v := strings.TrimSpace(input.CategoryIDs); v != "" {
		form["category_ids"] = v
	}
	if v := strings.TrimSpace(input.Keywords); v != "" {
		form["keywords"] = v
	}
	if v := strings.TrimSpace(input.MaxSalePrice); v != "" {
		form["max_sale_price"] = v
	}
	if v := strings.TrimSpace(input.MinSalePrice); v != "" {
		form["min_sale_price"] = v
	}
	if v := strings.TrimSpace(input.PageNo); v != "" {
		form["page_no"] = v
	}
	if v := strings.TrimSpace(input.PageSize); v != "" {
		form["page_size"] = v
	}
	if v := strings.TrimSpace(input.Sort); v != "" {
		form["sort"] = v
	}
	if v := strings.TrimSpace(input.TrackingID); v != "" {
		form["tracking_id"] = v
	}
	if v := strings.TrimSpace(input.ShipToCountry); v != "" {
		form["ship_to_country"] = v
	}

	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName:     "aliexpress.affiliate.product.query",
		form:        form,
		topProtocol: true,
	})
	if err != nil {
		return nil, err
	}

	var payload productQueryEnvelope
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode product query response: %w", err)
	}
	return &ProductQueryResult{
		CurrentPageNo:      payload.RespResult.Result.CurrentPageNo,
		CurrentRecordCount: payload.RespResult.Result.CurrentRecordCount,
		TotalPageNo:        payload.RespResult.Result.TotalPageNo,
		TotalRecordCount:   payload.RespResult.Result.TotalRecordCount,
		Products:           payload.RespResult.Result.Products.Product,
	}, nil
}

func (c *HTTPClient) GetAffiliateHotProducts(ctx context.Context, input ProductQueryInput) (*ProductQueryResult, error) {
	form := map[string]string{
		"target_currency": defaultString(input.TargetCurrency, "USD"),
		"target_language": defaultString(input.TargetLanguage, "EN"),
		"access_token":    strings.TrimSpace(input.AccessToken),
	}
	if v := strings.TrimSpace(input.CategoryIDs); v != "" {
		form["category_ids"] = v
	}
	if v := strings.TrimSpace(input.Keywords); v != "" {
		form["keywords"] = v
	}
	if v := strings.TrimSpace(input.MaxSalePrice); v != "" {
		form["max_sale_price"] = v
	}
	if v := strings.TrimSpace(input.MinSalePrice); v != "" {
		form["min_sale_price"] = v
	}
	if v := strings.TrimSpace(input.PageNo); v != "" {
		form["page_no"] = v
	}
	if v := strings.TrimSpace(input.PageSize); v != "" {
		form["page_size"] = v
	}
	if v := strings.TrimSpace(input.Sort); v != "" {
		form["sort"] = v
	}
	if v := strings.TrimSpace(input.TrackingID); v != "" {
		form["tracking_id"] = v
	}
	if v := strings.TrimSpace(input.ShipToCountry); v != "" {
		form["ship_to_country"] = v
	}

	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName:     "aliexpress.affiliate.hotproduct.query",
		form:        form,
		topProtocol: true,
	})
	if err != nil {
		return nil, err
	}

	var payload productQueryEnvelope
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode hot product query response: %w", err)
	}
	return &ProductQueryResult{
		CurrentPageNo:      payload.RespResult.Result.CurrentPageNo,
		CurrentRecordCount: payload.RespResult.Result.CurrentRecordCount,
		TotalPageNo:        payload.RespResult.Result.TotalPageNo,
		TotalRecordCount:   payload.RespResult.Result.TotalRecordCount,
		Products:           payload.RespResult.Result.Products.Product,
	}, nil
}

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
			"country":         strings.TrimSpace(input.Country),
			"tracking_id":     strings.TrimSpace(input.TrackingID),
			"access_token":    strings.TrimSpace(input.AccessToken),
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
			"access_token":      strings.TrimSpace(input.AccessToken),
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
