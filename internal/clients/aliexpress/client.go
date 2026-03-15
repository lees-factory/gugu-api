package aliexpress

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL   = "https://api-sg.aliexpress.com"
	defaultPartnerID = "gugu-api"
	signMethodSHA256 = "sha256"
)

type Config struct {
	BaseURL     string
	AppKey      string
	AppSecret   string
	CallbackURL string
	PartnerID   string
	HTTPClient  *http.Client
}

type RemoteError struct {
	Code      string
	Message   string
	RequestID string
}

func (e *RemoteError) Error() string {
	if e == nil {
		return ""
	}
	if e.Code == "" {
		return e.Message
	}
	return fmt.Sprintf("aliexpress error: code=%s message=%s", e.Code, e.Message)
}

type TokenExchangeInput struct {
	Code string
}

type RefreshTokenInput struct {
	RefreshToken string
}

type ProductLookupInput struct {
	ProductID      string
	TargetCurrency string
	TargetLanguage string
	Country        string
	TrackingID     string
	Fields         []string
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

type TokenSet struct {
	AccessToken           string
	RefreshToken          string
	ExpiresIn             int64
	RefreshExpiresIn      int64
	ExpireTime            int64
	RefreshTokenValidTime int64
	HavanaID              string
	UserID                string
	SellerID              string
	UserNick              string
	Account               string
	Locale                string
	AccountPlatform       string
	SP                    string
	RequestID             string
	Code                  string
}

type Client interface {
	BuildAuthorizationURL() (string, error)
	ExchangeCode(ctx context.Context, input TokenExchangeInput) (*TokenSet, error)
	RefreshAccessToken(ctx context.Context, input RefreshTokenInput) (*TokenSet, error)
	GetProductSnapshot(ctx context.Context, input ProductLookupInput) (*ProductSnapshot, error)
}

type HTTPClient struct {
	baseURL     string
	appKey      string
	appSecret   string
	callbackURL string
	partnerID   string
	httpClient  *http.Client
	now         func() time.Time
}

func NewHTTPClient(cfg Config) (*HTTPClient, error) {
	if strings.TrimSpace(cfg.AppKey) == "" {
		return nil, fmt.Errorf("app key is required")
	}
	if strings.TrimSpace(cfg.AppSecret) == "" {
		return nil, fmt.Errorf("app secret is required")
	}
	if strings.TrimSpace(cfg.CallbackURL) == "" {
		return nil, fmt.Errorf("callback url is required")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	partnerID := strings.TrimSpace(cfg.PartnerID)
	if partnerID == "" {
		partnerID = defaultPartnerID
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}

	return &HTTPClient{
		baseURL:     baseURL,
		appKey:      strings.TrimSpace(cfg.AppKey),
		appSecret:   strings.TrimSpace(cfg.AppSecret),
		callbackURL: strings.TrimSpace(cfg.CallbackURL),
		partnerID:   partnerID,
		httpClient:  httpClient,
		now:         time.Now,
	}, nil
}

func (c *HTTPClient) BuildAuthorizationURL() (string, error) {
	authURL, err := url.Parse(c.baseURL + "/oauth/authorize")
	if err != nil {
		return "", fmt.Errorf("parse authorization url: %w", err)
	}

	query := authURL.Query()
	query.Set("response_type", "code")
	query.Set("force_auth", "true")
	query.Set("redirect_uri", c.callbackURL)
	query.Set("client_id", c.appKey)
	authURL.RawQuery = query.Encode()

	return authURL.String(), nil
}

func (c *HTTPClient) ExchangeCode(ctx context.Context, input TokenExchangeInput) (*TokenSet, error) {
	code := strings.TrimSpace(input.Code)
	if code == "" {
		return nil, fmt.Errorf("code is required")
	}

	return c.executeTokenRequest(ctx, "/auth/token/create", map[string]string{
		"code": code,
	})
}

func (c *HTTPClient) RefreshAccessToken(ctx context.Context, input RefreshTokenInput) (*TokenSet, error) {
	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	return c.executeTokenRequest(ctx, "/auth/token/refresh", map[string]string{
		"refresh_token": refreshToken,
	})
}

func (c *HTTPClient) GetProductSnapshot(ctx context.Context, input ProductLookupInput) (*ProductSnapshot, error) {
	productID := strings.TrimSpace(input.ProductID)
	if productID == "" {
		return nil, fmt.Errorf("product id is required")
	}

	fields := input.Fields
	if len(fields) == 0 {
		fields = []string{
			"product_title",
			"target_sale_price",
			"target_sale_price_currency",
			"product_main_image_url",
			"promotion_link",
			"product_detail_url",
			"target_original_price",
			"shop_name",
		}
	}

	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName: "aliexpress.affiliate.productdetail.get",
		form: map[string]string{
			"fields":          strings.Join(fields, ","),
			"product_ids":     productID,
			"target_currency": defaultString(input.TargetCurrency, "USD"),
			"target_language": defaultString(input.TargetLanguage, "EN"),
			"country":         defaultString(input.Country, "US"),
			"tracking_id":     strings.TrimSpace(input.TrackingID),
		},
	})
	if err != nil {
		return nil, err
	}

	var payload productDetailEnvelope
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode product detail response: %w", err)
	}

	if payload.RespResult.RespCode != 20010000 {
		return nil, fmt.Errorf("aliexpress product detail failed: code=%d msg=%s", payload.RespResult.RespCode, payload.RespResult.RespMsg)
	}
	if len(payload.RespResult.Result.Products) == 0 {
		return nil, fmt.Errorf("aliexpress product detail returned no products")
	}

	product := payload.RespResult.Result.Products[0]
	price := firstNonEmpty(product.TargetSalePrice, product.SalePrice, product.TargetAppSalePrice, product.AppSalePrice)
	currency := firstNonEmpty(product.TargetSalePriceCurrency, product.SalePriceCurrency, product.TargetAppSalePriceCurrency, product.AppSalePriceCurrency)

	return &ProductSnapshot{
		ProductID:      strconv.FormatInt(product.ProductID, 10),
		Title:          product.ProductTitle,
		Price:          price,
		Currency:       currency,
		MainImageURL:   product.ProductMainImageURL,
		ProductURL:     product.ProductDetailURL,
		PromotionLink:  product.PromotionLink,
		OriginalPrice:  firstNonEmpty(product.TargetOriginalPrice, product.OriginalPrice),
		StoreName:      product.ShopName,
		TrackingIDUsed: strings.TrimSpace(input.TrackingID),
	}, nil
}

type signedRequest struct {
	apiName string
	form    map[string]string
}

func (c *HTTPClient) executeTokenRequest(ctx context.Context, apiName string, form map[string]string) (*TokenSet, error) {
	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName: apiName,
		form:    form,
	})
	if err != nil {
		return nil, err
	}

	var payload tokenResponse
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	if payload.AccessToken == "" {
		return nil, fmt.Errorf("aliexpress token response missing access token")
	}

	return &TokenSet{
		AccessToken:           payload.AccessToken,
		RefreshToken:          payload.RefreshToken,
		ExpiresIn:             payload.ExpiresIn,
		RefreshExpiresIn:      payload.RefreshExpiresIn,
		ExpireTime:            payload.ExpireTime,
		RefreshTokenValidTime: payload.RefreshTokenValidTime,
		HavanaID:              payload.HavanaID,
		UserID:                payload.UserID,
		SellerID:              payload.SellerID,
		UserNick:              payload.UserNick,
		Account:               payload.Account,
		Locale:                payload.Locale,
		AccountPlatform:       payload.AccountPlatform,
		SP:                    payload.SP,
		RequestID:             payload.RequestID,
		Code:                  payload.Code,
	}, nil
}

func (c *HTTPClient) executeFormRequest(ctx context.Context, request signedRequest) ([]byte, error) {
	form := make(url.Values)
	for key, value := range c.withoutEmptyValues(request.form) {
		form.Set(key, value)
	}

	query := url.Values{}
	query.Set("app_key", c.appKey)
	query.Set("sign_method", signMethodSHA256)
	query.Set("timestamp", strconv.FormatInt(c.now().UnixMilli(), 10))
	query.Set("partner_id", c.partnerID)
	query.Set("sign", c.sign(request.apiName, mergeValues(query, form)))

	endpoint := c.requestURL(request.apiName)
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"?"+query.Encode(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	httpRequest.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("request aliexpress api: %w", err)
	}
	defer response.Body.Close()

	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode raw response: %w", err)
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("re-encode raw response: %w", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("aliexpress api http %d: %s", response.StatusCode, string(raw))
	}

	if errResp, ok := payload["error_response"]; ok {
		errorRaw, _ := json.Marshal(errResp)
		var remoteErr remoteErrorEnvelope
		if err := json.Unmarshal(errorRaw, &remoteErr); err == nil {
			return nil, &RemoteError{
				Code:      remoteErr.Code,
				Message:   remoteErr.Message,
				RequestID: remoteErr.RequestID,
			}
		}
		return nil, fmt.Errorf("aliexpress api error: %s", string(errorRaw))
	}

	return raw, nil
}

func (c *HTTPClient) requestURL(apiName string) string {
	apiName = strings.TrimSpace(apiName)
	if strings.HasPrefix(apiName, "/") {
		return c.baseURL + "/rest" + apiName
	}

	return c.baseURL + "/rest/" + apiName
}

func (c *HTTPClient) sign(apiName string, values url.Values) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteString(apiName)
	for _, key := range keys {
		value := values.Get(key)
		if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			continue
		}
		builder.WriteString(key)
		builder.WriteString(value)
	}

	mac := hmac.New(sha256.New, []byte(c.appSecret))
	_, _ = mac.Write([]byte(builder.String()))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}

func (c *HTTPClient) withoutEmptyValues(values map[string]string) map[string]string {
	cleaned := make(map[string]string, len(values))
	for key, value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		cleaned[key] = strings.TrimSpace(value)
	}
	return cleaned
}

func mergeValues(values ...url.Values) url.Values {
	merged := make(url.Values)
	for _, current := range values {
		for key, currentValues := range current {
			for _, value := range currentValues {
				merged.Add(key, value)
			}
		}
	}
	return merged
}

func defaultString(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}

	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}

	return ""
}

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
	CurrentRecordCount int64                 `json:"current_record_count"`
	Products           []productDetailRecord `json:"products"`
}

type productDetailRecord struct {
	AppSalePrice               string `json:"app_sale_price"`
	AppSalePriceCurrency       string `json:"app_sale_price_currency"`
	OriginalPrice              string `json:"original_price"`
	ProductDetailURL           string `json:"product_detail_url"`
	ProductID                  int64  `json:"product_id"`
	ProductMainImageURL        string `json:"product_main_image_url"`
	ProductTitle               string `json:"product_title"`
	PromotionLink              string `json:"promotion_link"`
	SalePrice                  string `json:"sale_price"`
	SalePriceCurrency          string `json:"sale_price_currency"`
	ShopName                   string `json:"shop_name"`
	TargetAppSalePrice         string `json:"target_app_sale_price"`
	TargetAppSalePriceCurrency string `json:"target_app_sale_price_currency"`
	TargetOriginalPrice        string `json:"target_original_price"`
	TargetSalePrice            string `json:"target_sale_price"`
	TargetSalePriceCurrency    string `json:"target_sale_price_currency"`
}
