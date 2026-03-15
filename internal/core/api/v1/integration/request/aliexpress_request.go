package request

type AliExpressAuthorizationURL struct {
	UserID string `json:"user_id"`
}

type AliExpressExchangeCode struct {
	UserID string `json:"user_id"`
	Code   string `json:"code"`
}
