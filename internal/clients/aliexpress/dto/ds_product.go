package dto

type DSProductInput struct {
	ProductID             string
	ShipToCountry         string
	TargetCurrency        string
	TargetLanguage        string
	RemovePersonalBenefit string
	AccessToken           string
}

type DSProductResult struct {
	BaseInfo      DSItemBaseInfo     `json:"ae_item_base_info_dto"`
	SKUs          []DSItemSKUInfo    `json:"ae_item_sku_info_dtos"`
	Properties    []DSItemProperty   `json:"ae_item_properties"`
	Multimedia    DSMultimediaInfo   `json:"ae_multimedia_info_dto"`
	Logistics     DSLogisticsInfo    `json:"logistics_info_dto"`
	PackageInfo   DSPackageInfo      `json:"package_info_dto"`
	StoreInfo     DSStoreInfo        `json:"ae_store_info"`
	HasWholesale  bool               `json:"has_whole_sale"`
	RspCode       int64              `json:"rsp_code"`
	RspMsg        string             `json:"rsp_msg"`
}

type DSItemBaseInfo struct {
	ProductID          int64  `json:"product_id"`
	Subject            string `json:"subject"`
	ProductStatusType  string `json:"product_status_type"`
	CategoryID         int64  `json:"category_id"`
	CategorySequence   string `json:"category_sequence"`
	CurrencyCode       string `json:"currency_code"`
	AvgEvaluationRating string `json:"avg_evaluation_rating"`
	EvaluationCount    string `json:"evaluation_count"`
	SalesCount         string `json:"sales_count"`
	Detail             string `json:"detail"`
	MobileDetail       string `json:"mobile_detail"`
	GmtCreate          string `json:"gmt_create"`
	GmtModified        string `json:"gmt_modified"`
	OwnerMemberSeqLong int64  `json:"owner_member_seq_long"`
	SeparatedListing   bool   `json:"separated_listing"`
	SLProduct          bool   `json:"sl_product"`
	SLRelatedProductID int64  `json:"sl_related_product_id"`
}

type DSItemSKUInfo struct {
	SKUID                string              `json:"sku_id"`
	ID                   string              `json:"id"`
	SKUAttr              string              `json:"sku_attr"`
	SKUPrice             string              `json:"sku_price"`
	OfferSalePrice       string              `json:"offer_sale_price"`
	OfferBulkSalePrice   string              `json:"offer_bulk_sale_price"`
	ChannelDiscountPrice string              `json:"channel_discount_price"`
	SKUAvailableStock    int64               `json:"sku_available_stock"`
	SKUBulkOrder         int64               `json:"sku_bulk_order"`
	CurrencyCode         string              `json:"currency_code"`
	Barcode              string              `json:"barcode"`
	EANCode              string              `json:"ean_code"`
	PriceIncludeTax      bool                `json:"price_include_tax"`
	TaxCurrencyCode      string              `json:"tax_currency_code"`
	TaxAmount            string              `json:"tax_amount"`
	EstimatedImportCharges string            `json:"estimated_import_charges"`
	BuyAmountLimitByPromo  string            `json:"buy_amount_limit_set_by_promotion"`
	LimitStrategy        string              `json:"limit_strategy"`
	Properties           []DSSKUPropertyDTO  `json:"ae_sku_property_dtos"`
	WholesalePriceTiers  []DSWholesaleTier   `json:"wholesale_price_tiers"`
}

type DSSKUPropertyDTO struct {
	SKUPropertyID               int64  `json:"sku_property_id"`
	SKUPropertyName             string `json:"sku_property_name"`
	PropertyValueID             int64  `json:"property_value_id"`
	SKUPropertyValue            string `json:"sku_property_value"`
	PropertyValueDefinitionName string `json:"property_value_definition_name"`
	SKUImage                    string `json:"sku_image"`
}

type DSWholesaleTier struct {
	MinQuantity    string `json:"min_quantity"`
	Discount       string `json:"discount"`
	WholesalePrice string `json:"wholesale_price"`
}

type DSItemProperty struct {
	AttrNameID      int64  `json:"attr_name_id"`
	AttrName        string `json:"attr_name"`
	AttrValueID     int64  `json:"attr_value_id"`
	AttrValue       string `json:"attr_value"`
	AttrValueUnit   string `json:"attr_value_unit"`
	AttrValueStart  string `json:"attr_value_start"`
	AttrValueEnd    string `json:"attr_value_end"`
}

type DSMultimediaInfo struct {
	ImageURLs string       `json:"image_urls"`
	Videos    []DSVideoDTO `json:"ae_video_dtos"`
}

type DSVideoDTO struct {
	MediaID    int64  `json:"media_id"`
	MediaURL   string `json:"media_url"`
	PosterURL  string `json:"poster_url"`
	MediaType  string `json:"media_type"`
	MediaStatus string `json:"media_status"`
	AliMemberID int64  `json:"ali_member_id"`
}

type DSLogisticsInfo struct {
	ShipToCountry string `json:"ship_to_country"`
	DeliveryTime  int64  `json:"delivery_time"`
}

type DSPackageInfo struct {
	PackageLength int64  `json:"package_length"`
	PackageWidth  int64  `json:"package_width"`
	PackageHeight int64  `json:"package_height"`
	GrossWeight   string `json:"gross_weight"`
	BaseUnit      int64  `json:"base_unit"`
	ProductUnit   int64  `json:"product_unit"`
	PackageType   bool   `json:"package_type"`
}

type DSStoreInfo struct {
	StoreID                int64  `json:"store_id"`
	StoreName              string `json:"store_name"`
	StoreCountryCode       string `json:"store_country_code"`
	ItemAsDescribedRating  string `json:"item_as_described_rating"`
	CommunicationRating    string `json:"communication_rating"`
	ShippingSpeedRating    string `json:"shipping_speed_rating"`
}
