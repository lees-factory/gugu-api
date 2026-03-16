package request

type AddTrackedItem struct {
	UserID            string `json:"user_id"`
	OriginalURL       string `json:"original_url"`
	ProviderCommerce  string `json:"provider_commerce"`
	ExternalProductID string `json:"external_product_id"`
}
