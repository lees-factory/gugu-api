package request

type LoadHotProductsRequest struct {
	CategoryIDs    string `json:"category_ids"`
	Keywords       string `json:"keywords"`
	PageNo         string `json:"page_no"`
	PageSize       string `json:"page_size"`
	Sort           string `json:"sort"`
	MinSalePrice   string `json:"min_sale_price"`
	MaxSalePrice   string `json:"max_sale_price"`
	ShipToCountry  string `json:"ship_to_country"`
	TargetCurrency string `json:"target_currency"`
	TargetLanguage string `json:"target_language"`
}
