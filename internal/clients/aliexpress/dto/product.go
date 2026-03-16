package dto

type ProductLookupInput struct {
	ProductID      string
	TargetCurrency string
	TargetLanguage string
	Country        string
	TrackingID     string
	Fields         []string
}

type ProductDetailInput struct {
	ProductIDs     []string
	TargetCurrency string
	TargetLanguage string
	Country        string
	TrackingID     string
	Fields         []string
}

type ProductDetailResult struct {
	CurrentRecordCount int64
	Products           []AffiliateProduct
}

type ProductSKUDetailInput struct {
	ProductID       string
	ShipToCountry   string
	TargetCurrency  string
	TargetLanguage  string
	NeedDeliverInfo string
	SKUIDs          []string
}

type ProductSKUDetailResult struct {
	Code     int64
	Success  bool
	ItemInfo AffiliateSKUItemInfo
	SKUInfos []AffiliateSKUInfo
}

type ProductSnapshot struct {
	ProductID      string
	Title          string
	Price          string
	Currency       string
	MainImageURL   string
	ProductURL     string
	PromotionLink  string
	OriginalPrice  string
	StoreName      string
	TrackingIDUsed string
}

type AffiliateProduct struct {
	SKUID                        int64         `json:"sku_id"`
	TaxRate                      string        `json:"tax_rate"`
	AppSalePrice                 string        `json:"app_sale_price"`
	AppSalePriceCurrency         string        `json:"app_sale_price_currency"`
	CommissionRate               string        `json:"commission_rate"`
	EANCode                      string        `json:"ean_code"`
	Discount                     string        `json:"discount"`
	EvaluateRate                 string        `json:"evaluate_rate"`
	FirstLevelCategoryID         int64         `json:"first_level_category_id"`
	FirstLevelCategoryName       string        `json:"first_level_category_name"`
	LastestVolume                int64         `json:"lastest_volume"`
	HotProductCommissionRate     string        `json:"hot_product_commission_rate"`
	OriginalPrice                string        `json:"original_price"`
	OriginalPriceCurrency        string        `json:"original_price_currency"`
	PlatformProductType          string        `json:"platform_product_type"`
	ProductDetailURL             string        `json:"product_detail_url"`
	ProductID                    int64         `json:"product_id"`
	ProductMainImageURL          string        `json:"product_main_image_url"`
	ProductSmallImageURLs        []string      `json:"product_small_image_urls"`
	ProductTitle                 string        `json:"product_title"`
	ProductVideoURL              string        `json:"product_video_url"`
	PromotionLink                string        `json:"promotion_link"`
	SalePrice                    string        `json:"sale_price"`
	SalePriceCurrency            string        `json:"sale_price_currency"`
	SecondLevelCategoryID        int64         `json:"second_level_category_id"`
	SecondLevelCategoryName      string        `json:"second_level_category_name"`
	ShopName                     string        `json:"shop_name"`
	ShopID                       int64         `json:"shop_id"`
	ShopURL                      string        `json:"shop_url"`
	TargetAppSalePrice           string        `json:"target_app_sale_price"`
	TargetAppSalePriceCurrency   string        `json:"target_app_sale_price_currency"`
	TargetOriginalPrice          string        `json:"target_original_price"`
	TargetOriginalPriceCurrency  string        `json:"target_original_price_currency"`
	TargetSalePrice              string        `json:"target_sale_price"`
	TargetSalePriceCurrency      string        `json:"target_sale_price_currency"`
	RelevantMarketCommissionRate string        `json:"relevant_market_commission_rate"`
	PromoCodeInfo                PromoCodeInfo `json:"promo_code_info"`
}

type PromoCodeInfo struct {
	PromoCode              string `json:"promo_code"`
	CodeCampaignType       string `json:"code_campaigntype"`
	CodeValue              string `json:"code_value"`
	CodeAvailableTimeStart string `json:"code_availabletime_start"`
	CodeAvailableTimeEnd   string `json:"code_availabletime_end"`
	CodeMiniSpend          string `json:"code_mini_spend"`
	CodeQuantity           string `json:"code_quantity"`
	CodePromotionURL       string `json:"code_promotionurl"`
}

type AffiliateSKUItemInfo struct {
	OrderNumber           string   `json:"order_number"`
	ImageWhite            string   `json:"image_white"`
	ProductID             string   `json:"product_id"`
	DisplayCategoryNameL1 string   `json:"display_category_name_l1"`
	DisplayCategoryNameL2 string   `json:"display_category_name_l2"`
	DisplayCategoryNameL3 string   `json:"display_category_name_l3"`
	DisplayCategoryNameL4 string   `json:"display_category_name_l4"`
	ProductScore          string   `json:"product_score"`
	Title                 string   `json:"title"`
	OriginalLink          string   `json:"original_link"`
	ProductCategory       string   `json:"product_category"`
	ImageLink             string   `json:"image_link"`
	EnTitle               string   `json:"en_title"`
	AdditionalImageLinks  []string `json:"additional_image_links"`
	ReviewNumber          string   `json:"review_number"`
	Brand                 string   `json:"brand"`
	AgeGroup              string   `json:"age_group"`
	Gender                string   `json:"gender"`
	Condition             string   `json:"condition"`
	StoreName             string   `json:"store_name"`
}

type AffiliateSKUInfo struct {
	DiscountRate     string `json:"discount_rate"`
	Link             string `json:"link"`
	ShippingFees     string `json:"shipping_fees"`
	Color            string `json:"color"`
	MaxDeliveryDays  string `json:"max_delivery_days"`
	PriceWithTax     string `json:"price_with_tax"`
	MinDeliveryDays  string `json:"min_delivery_days"`
	ShipFromCountry  string `json:"ship_from_country"`
	SalePriceWithTax string `json:"sale_price_with_tax"`
	TaxRate          string `json:"tax_rate"`
	SKUImageLink     string `json:"sku_image_link"`
	Size             string `json:"size"`
	DeliveryDays     string `json:"delivery_days"`
	Currency         string `json:"currency"`
	SKUID            int64  `json:"sku_id"`
	EANCode          string `json:"ean_code"`
	SKUProperties    string `json:"sku_properties"`
}
