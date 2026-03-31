package aliexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (c *HTTPClient) GetDSProduct(ctx context.Context, input DSProductInput) (*DSProductResult, error) {
	productID := strings.TrimSpace(input.ProductID)
	if productID == "" {
		return nil, fmt.Errorf("product id is required")
	}
	shipToCountry := strings.TrimSpace(input.ShipToCountry)
	if shipToCountry == "" {
		return nil, fmt.Errorf("ship to country is required")
	}

	form := map[string]string{
		"product_id":      productID,
		"ship_to_country": shipToCountry,
		"target_currency": defaultString(input.TargetCurrency, "USD"),
		"target_language": defaultString(input.TargetLanguage, "en"),
		"access_token":    strings.TrimSpace(input.AccessToken),
	}
	if strings.TrimSpace(input.RemovePersonalBenefit) != "" {
		form["remove_personal_benefit"] = input.RemovePersonalBenefit
	}

	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName:     "aliexpress.ds.product.get",
		form:        form,
		topProtocol: true,
	})
	if err != nil {
		return nil, err
	}

	var payload dsProductEnvelope
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode ds product response: %w", err)
	}

	inner := &payload.Result
	if inner.RspCode != 0 && inner.RspCode != 200 {
		return nil, &RemoteError{
			Code:    fmt.Sprintf("%d", inner.RspCode),
			Message: inner.RspMsg,
		}
	}

	return toDSProductResult(inner), nil
}

func toDSProductResult(r *dsProductResult) *DSProductResult {
	skus := make([]DSItemSKUInfo, 0, len(r.SKUs.SKUInfos))
	for _, s := range r.SKUs.SKUInfos {
		props := make([]DSSKUPropertyDTO, 0, len(s.Properties.Properties))
		for _, p := range s.Properties.Properties {
			props = append(props, DSSKUPropertyDTO{
				SKUPropertyID:               p.SKUPropertyID,
				SKUPropertyName:             p.SKUPropertyName,
				PropertyValueID:             p.PropertyValueID,
				SKUPropertyValue:            p.SKUPropertyValue,
				PropertyValueDefinitionName: p.PropertyValueDefinitionName,
				SKUImage:                    p.SKUImage,
			})
		}
		tiers := make([]DSWholesaleTier, 0, len(s.WholesalePriceTiers.Tiers))
		for _, t := range s.WholesalePriceTiers.Tiers {
			tiers = append(tiers, DSWholesaleTier{
				MinQuantity:    t.MinQuantity,
				Discount:       t.Discount,
				WholesalePrice: t.WholesalePrice,
			})
		}
		skus = append(skus, DSItemSKUInfo{
			SKUID:                  s.SKUID,
			ID:                     s.ID,
			SKUAttr:                s.SKUAttr,
			SKUPrice:               s.SKUPrice,
			OfferSalePrice:         s.OfferSalePrice,
			OfferBulkSalePrice:     s.OfferBulkSalePrice,
			ChannelDiscountPrice:   s.ChannelDiscountPrice,
			SKUAvailableStock:      s.SKUAvailableStock,
			SKUBulkOrder:           s.SKUBulkOrder,
			CurrencyCode:           s.CurrencyCode,
			Barcode:                s.Barcode,
			EANCode:                s.EANCode,
			PriceIncludeTax:        s.PriceIncludeTax,
			TaxCurrencyCode:        s.TaxCurrencyCode,
			TaxAmount:              s.TaxAmount,
			EstimatedImportCharges: s.EstimatedImportCharges,
			BuyAmountLimitByPromo:  s.BuyAmountLimitByPromo,
			LimitStrategy:          s.LimitStrategy,
			Properties:             props,
			WholesalePriceTiers:    tiers,
		})
	}

	properties := make([]DSItemProperty, 0, len(r.Properties.Properties))
	for _, p := range r.Properties.Properties {
		properties = append(properties, DSItemProperty{
			AttrNameID:     p.AttrNameID,
			AttrName:       p.AttrName,
			AttrValueID:    p.AttrValueID,
			AttrValue:      p.AttrValue,
			AttrValueUnit:  p.AttrValueUnit,
			AttrValueStart: p.AttrValueStart,
			AttrValueEnd:   p.AttrValueEnd,
		})
	}

	videos := make([]DSVideoDTO, 0, len(r.Multimedia.Videos.Videos))
	for _, v := range r.Multimedia.Videos.Videos {
		videos = append(videos, DSVideoDTO{
			MediaID:     v.MediaID,
			MediaURL:    v.MediaURL,
			PosterURL:   v.PosterURL,
			MediaType:   v.MediaType,
			MediaStatus: v.MediaStatus,
			AliMemberID: v.AliMemberID,
		})
	}

	return &DSProductResult{
		BaseInfo: DSItemBaseInfo{
			ProductID:           r.BaseInfo.ProductID,
			Subject:             r.BaseInfo.Subject,
			ProductStatusType:   r.BaseInfo.ProductStatusType,
			CategoryID:          r.BaseInfo.CategoryID,
			CategorySequence:    r.BaseInfo.CategorySequence,
			CurrencyCode:        r.BaseInfo.CurrencyCode,
			AvgEvaluationRating: r.BaseInfo.AvgEvaluationRating,
			EvaluationCount:     r.BaseInfo.EvaluationCount,
			SalesCount:          r.BaseInfo.SalesCount,
			Detail:              r.BaseInfo.Detail,
			MobileDetail:        r.BaseInfo.MobileDetail,
			GmtCreate:           r.BaseInfo.GmtCreate,
			GmtModified:         r.BaseInfo.GmtModified,
			OwnerMemberSeqLong:  r.BaseInfo.OwnerMemberSeqLong,
			SeparatedListing:    r.BaseInfo.SeparatedListing,
			SLProduct:           r.BaseInfo.SLProduct,
			SLRelatedProductID:  r.BaseInfo.SLRelatedProductID,
		},
		SKUs:       skus,
		Properties: properties,
		Multimedia: DSMultimediaInfo{
			ImageURLs: r.Multimedia.ImageURLs,
			Videos:    videos,
		},
		Logistics: DSLogisticsInfo{
			ShipToCountry: r.Logistics.ShipToCountry,
			DeliveryTime:  r.Logistics.DeliveryTime,
		},
		PackageInfo: DSPackageInfo{
			PackageLength: r.PackageInfo.PackageLength,
			PackageWidth:  r.PackageInfo.PackageWidth,
			PackageHeight: r.PackageInfo.PackageHeight,
			GrossWeight:   r.PackageInfo.GrossWeight,
			BaseUnit:      r.PackageInfo.BaseUnit,
			ProductUnit:   r.PackageInfo.ProductUnit,
			PackageType:   r.PackageInfo.PackageType,
		},
		StoreInfo: DSStoreInfo{
			StoreID:               r.StoreInfo.StoreID,
			StoreName:             r.StoreInfo.StoreName,
			StoreCountryCode:      r.StoreInfo.StoreCountryCode,
			ItemAsDescribedRating: r.StoreInfo.ItemAsDescribedRating,
			CommunicationRating:   r.StoreInfo.CommunicationRating,
			ShippingSpeedRating:   r.StoreInfo.ShippingSpeedRating,
		},
		HasWholesale: r.HasWholesale,
		RspCode:      r.RspCode,
		RspMsg:       r.RspMsg,
	}
}
