package request

type SelectSKU struct {
	UserID string `json:"user_id"`
	SKUID  string `json:"sku_id"`
}
