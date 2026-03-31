package request

type AliExpressExchangeCode struct {
	Code    string `json:"code"`
	AppType string `json:"app_type"`
}

type AliExpressAuthorizeURL struct {
	AppType string `json:"app_type"`
}
