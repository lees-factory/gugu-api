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

type signedRequest struct {
	apiName     string
	form        map[string]string
	topProtocol bool
}

func (c *HTTPClient) executeFormRequest(ctx context.Context, request signedRequest) ([]byte, error) {
	if request.topProtocol {
		return c.executeTopRequest(ctx, request)
	}
	return c.executeLegacyRequest(ctx, request)
}

func (c *HTTPClient) executeTopRequest(ctx context.Context, request signedRequest) ([]byte, error) {
	allParams := map[string]string{
		"app_key":     c.appKey,
		"sign_method": signMethodSHA256,
		"timestamp":   strconv.FormatInt(c.now().UnixMilli(), 10),
		"method":      request.apiName,
		"v":           "2.0",
		"format":      "json",
		"partner_id":  c.partnerID,
	}
	for key, value := range c.withoutEmptyValues(request.form) {
		allParams[key] = value
	}
	allParams["sign"] = c.signTop(allParams)

	form := make(url.Values)
	for key, value := range allParams {
		form.Set(key, value)
	}

	endpoint := c.baseURL + "/sync"
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
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

	if response.StatusCode >= http.StatusBadRequest {
		raw, _ := json.Marshal(payload)
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

	responseKey := strings.ReplaceAll(request.apiName, ".", "_") + "_response"
	inner, ok := payload[responseKey]
	if !ok {
		raw, _ := json.Marshal(payload)
		return nil, fmt.Errorf("aliexpress api: missing response key %q: %s", responseKey, string(raw))
	}

	raw, err := json.Marshal(inner)
	if err != nil {
		return nil, fmt.Errorf("re-encode inner response: %w", err)
	}

	return raw, nil
}

func (c *HTTPClient) executeLegacyRequest(ctx context.Context, request signedRequest) ([]byte, error) {
	form := make(url.Values)
	for key, value := range c.withoutEmptyValues(request.form) {
		form.Set(key, value)
	}

	query := url.Values{}
	query.Set("app_key", c.appKey)
	query.Set("sign_method", signMethodSHA256)
	query.Set("timestamp", strconv.FormatInt(c.now().UnixMilli(), 10))
	query.Set("partner_id", c.partnerID)
	query.Set("sign", c.signLegacy(request.apiName, mergeValues(query, form)))

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

func (c *HTTPClient) signTop(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if key == "sign" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		value := params[key]
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

func (c *HTTPClient) signLegacy(apiName string, values url.Values) string {
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

func normalizeCSV(values []string) []string {
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		cleaned = append(cleaned, value)
	}
	return cleaned
}
