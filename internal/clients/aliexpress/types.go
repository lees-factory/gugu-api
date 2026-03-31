package aliexpress

type tokenResponse struct {
	RefreshTokenValidTime int64  `json:"refresh_token_valid_time"`
	ExpireTime            int64  `json:"expire_time"`
	HavanaID              string `json:"havana_id"`
	Locale                string `json:"locale"`
	UserNick              string `json:"user_nick"`
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	UserID                string `json:"user_id"`
	AccountPlatform       string `json:"account_platform"`
	RefreshExpiresIn      int64  `json:"refresh_expires_in"`
	ExpiresIn             int64  `json:"expires_in"`
	SP                    string `json:"sp"`
	SellerID              string `json:"seller_id"`
	Account               string `json:"account"`
	Code                  string `json:"code"`
	RequestID             string `json:"request_id"`
}

type remoteErrorEnvelope struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type categoryEnvelope struct {
	RespResult categoryResult `json:"resp_result"`
}

type categoryResult struct {
	RespCode int64               `json:"resp_code"`
	RespMsg  string              `json:"resp_msg"`
	Result   categoryResultBlock `json:"result"`
}

type categoryResultBlock struct {
	TotalResultCount int64               `json:"total_result_count"`
	Categories       categoryListWrapper `json:"categories"`
}

type categoryListWrapper struct {
	Category []AffiliateCategory `json:"category"`
}

type productQueryEnvelope struct {
	RespResult productQueryResult `json:"resp_result"`
}

type productQueryResult struct {
	RespCode int64                    `json:"resp_code"`
	RespMsg  string                   `json:"resp_msg"`
	Result   productQueryProductBlock `json:"result"`
}

type productQueryProductBlock struct {
	CurrentPageNo      int64                       `json:"current_page_no"`
	CurrentRecordCount int64                       `json:"current_record_count"`
	TotalPageNo        int64                       `json:"total_page_no"`
	TotalRecordCount   int64                       `json:"total_record_count"`
	Products           productDetailProductWrapper `json:"products"`
}

type productDetailEnvelope struct {
	RespResult productDetailResult `json:"resp_result"`
}

type productDetailResult struct {
	RespCode int64                     `json:"resp_code"`
	RespMsg  string                    `json:"resp_msg"`
	Result   productDetailProductBlock `json:"result"`
}

type productDetailProductBlock struct {
	CurrentRecordCount int64                       `json:"current_record_count"`
	Products           productDetailProductWrapper `json:"products"`
}

type productSKUDetailEnvelope struct {
	Result productSKUDetailOuterResult `json:"result"`
}

type productSKUDetailOuterResult struct {
	Result productSKUDetailInnerResult `json:"result"`
}

type productSKUDetailInnerResult struct {
	ItemInfo AffiliateSKUItemInfo  `json:"ae_item_info"`
	SKUInfos skuInfoTrafficWrapper `json:"ae_item_sku_info"`
	Code     int64                 `json:"code"`
	Success  bool                  `json:"success"`
}

type skuInfoTrafficWrapper struct {
	TrafficSKUInfoList []AffiliateSKUInfo `json:"traffic_sku_info_list"`
}

type productDetailProductWrapper struct {
	Product []AffiliateProduct `json:"product"`
}

// --- DS API envelope types ---

type dsProductEnvelope struct {
	Result dsProductResult `json:"result"`
}

type dsProductResult struct {
	BaseInfo      dsItemBaseInfo        `json:"ae_item_base_info_dto"`
	SKUs          dsItemSKUInfoWrapper  `json:"ae_item_sku_info_dtos"`
	Properties    dsItemPropertyWrapper `json:"ae_item_properties"`
	Multimedia    dsMultimediaInfo      `json:"ae_multimedia_info_dto"`
	Logistics     dsLogisticsInfo       `json:"logistics_info_dto"`
	PackageInfo   dsPackageInfo         `json:"package_info_dto"`
	StoreInfo     dsStoreInfo           `json:"ae_store_info"`
	HasWholesale  bool                  `json:"has_whole_sale"`
	RspCode       int64                 `json:"rsp_code"`
	RspMsg        string                `json:"rsp_msg"`
}

type dsItemBaseInfo struct {
	ProductID           int64  `json:"product_id"`
	Subject             string `json:"subject"`
	ProductStatusType   string `json:"product_status_type"`
	CategoryID          int64  `json:"category_id"`
	CategorySequence    string `json:"category_sequence"`
	CurrencyCode        string `json:"currency_code"`
	AvgEvaluationRating string `json:"avg_evaluation_rating"`
	EvaluationCount     string `json:"evaluation_count"`
	SalesCount          string `json:"sales_count"`
	Detail              string `json:"detail"`
	MobileDetail        string `json:"mobile_detail"`
	GmtCreate           string `json:"gmt_create"`
	GmtModified         string `json:"gmt_modified"`
	OwnerMemberSeqLong  int64  `json:"owner_member_seq_long"`
	SeparatedListing    bool   `json:"separated_listing"`
	SLProduct           bool   `json:"sl_product"`
	SLRelatedProductID  int64  `json:"sl_related_product_id"`
}

type dsItemSKUInfoWrapper struct {
	SKUInfos []dsItemSKUInfo `json:"ae_item_sku_info_d_t_o"`
}

type dsItemSKUInfo struct {
	SKUID                  string                   `json:"sku_id"`
	ID                     string                   `json:"id"`
	SKUAttr                string                   `json:"sku_attr"`
	SKUPrice               string                   `json:"sku_price"`
	OfferSalePrice         string                   `json:"offer_sale_price"`
	OfferBulkSalePrice     string                   `json:"offer_bulk_sale_price"`
	ChannelDiscountPrice   string                   `json:"channel_discount_price"`
	SKUAvailableStock      int64                    `json:"sku_available_stock"`
	SKUBulkOrder           int64                    `json:"sku_bulk_order"`
	CurrencyCode           string                   `json:"currency_code"`
	Barcode                string                   `json:"barcode"`
	EANCode                string                   `json:"ean_code"`
	PriceIncludeTax        bool                     `json:"price_include_tax"`
	TaxCurrencyCode        string                   `json:"tax_currency_code"`
	TaxAmount              string                   `json:"tax_amount"`
	EstimatedImportCharges string                   `json:"estimated_import_charges"`
	BuyAmountLimitByPromo  string                   `json:"buy_amount_limit_set_by_promotion"`
	LimitStrategy          string                   `json:"limit_strategy"`
	Properties             dsSKUPropertyDTOWrapper  `json:"ae_sku_property_dtos"`
	WholesalePriceTiers    dsWholesaleTierWrapper   `json:"wholesale_price_tiers"`
}

type dsSKUPropertyDTOWrapper struct {
	Properties []dsSKUPropertyDTO `json:"ae_sku_property_d_t_o"`
}

type dsSKUPropertyDTO struct {
	SKUPropertyID               int64  `json:"sku_property_id"`
	SKUPropertyName             string `json:"sku_property_name"`
	PropertyValueID             int64  `json:"property_value_id"`
	SKUPropertyValue            string `json:"sku_property_value"`
	PropertyValueDefinitionName string `json:"property_value_definition_name"`
	SKUImage                    string `json:"sku_image"`
}

type dsWholesaleTierWrapper struct {
	Tiers []dsWholesaleTier `json:"wholesale_price_tier"`
}

type dsWholesaleTier struct {
	MinQuantity    string `json:"min_quantity"`
	Discount       string `json:"discount"`
	WholesalePrice string `json:"wholesale_price"`
}

type dsItemPropertyWrapper struct {
	Properties []dsItemProperty `json:"ae_item_property"`
}

type dsItemProperty struct {
	AttrNameID     int64  `json:"attr_name_id"`
	AttrName       string `json:"attr_name"`
	AttrValueID    int64  `json:"attr_value_id"`
	AttrValue      string `json:"attr_value"`
	AttrValueUnit  string `json:"attr_value_unit"`
	AttrValueStart string `json:"attr_value_start"`
	AttrValueEnd   string `json:"attr_value_end"`
}

type dsMultimediaInfo struct {
	ImageURLs string          `json:"image_urls"`
	Videos    dsVideoWrapper  `json:"ae_video_dtos"`
}

type dsVideoWrapper struct {
	Videos []dsVideoDTO `json:"ae_video_d_t_o"`
}

type dsVideoDTO struct {
	MediaID     int64  `json:"media_id"`
	MediaURL    string `json:"media_url"`
	PosterURL   string `json:"poster_url"`
	MediaType   string `json:"media_type"`
	MediaStatus string `json:"media_status"`
	AliMemberID int64  `json:"ali_member_id"`
}

type dsLogisticsInfo struct {
	ShipToCountry string `json:"ship_to_country"`
	DeliveryTime  int64  `json:"delivery_time"`
}

type dsPackageInfo struct {
	PackageLength int64  `json:"package_length"`
	PackageWidth  int64  `json:"package_width"`
	PackageHeight int64  `json:"package_height"`
	GrossWeight   string `json:"gross_weight"`
	BaseUnit      int64  `json:"base_unit"`
	ProductUnit   int64  `json:"product_unit"`
	PackageType   bool   `json:"package_type"`
}

type dsStoreInfo struct {
	StoreID               int64  `json:"store_id"`
	StoreName             string `json:"store_name"`
	StoreCountryCode      string `json:"store_country_code"`
	ItemAsDescribedRating string `json:"item_as_described_rating"`
	CommunicationRating   string `json:"communication_rating"`
	ShippingSpeedRating   string `json:"shipping_speed_rating"`
}
