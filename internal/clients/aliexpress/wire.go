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

type productDetailEnvelope struct {
	RespResult productDetailResult `json:"resp_result"`
}

type productDetailResult struct {
	RespCode int64                     `json:"resp_code"`
	RespMsg  string                    `json:"resp_msg"`
	Result   productDetailProductBlock `json:"result"`
}

type productDetailProductBlock struct {
	CurrentRecordCount int64              `json:"current_record_count"`
	Products           []AffiliateProduct `json:"products"`
}

type productSKUDetailEnvelope struct {
	Result productSKUDetailOuterResult `json:"result"`
}

type productSKUDetailOuterResult struct {
	Result productSKUDetailInnerResult `json:"result"`
}

type productSKUDetailInnerResult struct {
	ItemInfo AffiliateSKUItemInfo `json:"ae_item_info"`
	SKUInfos []AffiliateSKUInfo   `json:"ae_item_sku_info"`
	Code     int64                `json:"code"`
	Success  bool                 `json:"success"`
}
