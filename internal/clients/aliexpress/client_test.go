package aliexpress

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestBuildAuthorizationURL(t *testing.T) {
	client := newTestClient(t, nil)

	got, err := client.BuildAuthorizationURL()
	if err != nil {
		t.Fatalf("BuildAuthorizationURL() error = %v", err)
	}

	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if parsed.Path != "/oauth/authorize" {
		t.Fatalf("path = %q, want /oauth/authorize", parsed.Path)
	}

	query := parsed.Query()
	if query.Get("client_id") != "528586" {
		t.Fatalf("client_id = %q, want 528586", query.Get("client_id"))
	}
	if query.Get("redirect_uri") != "https://googoo-client.vercel.app/callback" {
		t.Fatalf("redirect_uri = %q", query.Get("redirect_uri"))
	}
	if query.Get("response_type") != "code" {
		t.Fatalf("response_type = %q, want code", query.Get("response_type"))
	}
}

func TestExchangeCode(t *testing.T) {
	var capturedPath string
	var capturedQuery url.Values
	var capturedBody url.Values

	client := newTestClient(t, &Config{
		BaseURL:     "https://api-sg.aliexpress.com",
		AppKey:      "528586",
		AppSecret:   "secret",
		CallbackURL: "https://googoo-client.vercel.app/callback",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				capturedPath = r.URL.Path
				capturedQuery = r.URL.Query()
				capturedBody, err = url.ParseQuery(string(body))
				if err != nil {
					t.Fatalf("ParseQuery() error = %v", err)
				}

				return jsonResponse(`{
			"refresh_token_valid_time": 4804563378176,
			"expire_time": 4804563378176,
			"locale": "zh_CN",
			"user_nick": "test1234",
			"access_token": "access-token",
			"refresh_token": "refresh-token",
			"user_id": "200042362",
			"account_platform": "seller_center",
			"refresh_expires_in": 3153599922,
			"expires_in": 3153599922,
			"seller_id": "200042362",
			"account": "test1234@126.com",
			"code": "0",
			"request_id": "request-1"
		}`)
			}),
		},
	})

	tokens, err := client.ExchangeCode(context.Background(), TokenExchangeInput{Code: "auth-code"})
	if err != nil {
		t.Fatalf("ExchangeCode() error = %v", err)
	}

	if capturedPath != "/rest/auth/token/create" {
		t.Fatalf("path = %q, want /rest/auth/token/create", capturedPath)
	}
	if capturedBody.Get("code") != "auth-code" {
		t.Fatalf("body code = %q, want auth-code", capturedBody.Get("code"))
	}
	if capturedQuery.Get("app_key") != "528586" {
		t.Fatalf("app_key = %q, want 528586", capturedQuery.Get("app_key"))
	}
	if capturedQuery.Get("sign") == "" {
		t.Fatal("sign was empty")
	}
	if tokens.AccessToken != "access-token" {
		t.Fatalf("AccessToken = %q, want access-token", tokens.AccessToken)
	}
	if tokens.RefreshToken != "refresh-token" {
		t.Fatalf("RefreshToken = %q, want refresh-token", tokens.RefreshToken)
	}
}

func TestGetProductSnapshot(t *testing.T) {
	client := newTestClient(t, &Config{
		BaseURL:     "https://api-sg.aliexpress.com",
		AppKey:      "528586",
		AppSecret:   "secret",
		CallbackURL: "https://googoo-client.vercel.app/callback",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != "/sync" {
					t.Fatalf("path = %q, want /sync", r.URL.Path)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				values, err := url.ParseQuery(string(body))
				if err != nil {
					t.Fatalf("ParseQuery() error = %v", err)
				}
				if values.Get("product_ids") != "1005001234567890" {
					t.Fatalf("product_ids = %q", values.Get("product_ids"))
				}
				if values.Get("method") != "aliexpress.affiliate.productdetail.get" {
					t.Fatalf("method = %q", values.Get("method"))
				}
				if values.Get("v") != "2.0" {
					t.Fatalf("v = %q", values.Get("v"))
				}

				return jsonResponse(`{
			"aliexpress_affiliate_productdetail_get_response": {
				"resp_result": {
					"resp_code": 20010000,
					"resp_msg": "success",
					"result": {
						"current_record_count": 1,
						"products": {
							"product": [
								{
									"product_id": 1005001234567890,
									"product_title": "Keyboard",
									"target_sale_price": "89.99",
									"target_sale_price_currency": "USD",
									"product_main_image_url": "https://img.example/1.jpg",
									"product_detail_url": "https://www.aliexpress.com/item/1005001234567890.html",
									"promotion_link": "https://s.click.aliexpress.com/promo",
									"target_original_price": "99.99",
									"shop_name": "Gugu Store"
								}
							]
						}
					}
				}
			}
		}`)
			}),
		},
	})

	snapshot, err := client.GetProductSnapshot(context.Background(), ProductLookupInput{
		ProductID: "1005001234567890",
	})
	if err != nil {
		t.Fatalf("GetProductSnapshot() error = %v", err)
	}

	if snapshot.ProductID != "1005001234567890" {
		t.Fatalf("ProductID = %q", snapshot.ProductID)
	}
	if snapshot.Title != "Keyboard" {
		t.Fatalf("Title = %q", snapshot.Title)
	}
	if snapshot.Price != "89.99" {
		t.Fatalf("Price = %q", snapshot.Price)
	}
	if snapshot.Currency != "USD" {
		t.Fatalf("Currency = %q", snapshot.Currency)
	}
}

func TestGetAffiliateProductDetail(t *testing.T) {
	client := newTestClient(t, &Config{
		BaseURL:     "https://api-sg.aliexpress.com",
		AppKey:      "528586",
		AppSecret:   "secret",
		CallbackURL: "https://googoo-client.vercel.app/callback",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != "/sync" {
					t.Fatalf("path = %q, want /sync", r.URL.Path)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				values, err := url.ParseQuery(string(body))
				if err != nil {
					t.Fatalf("ParseQuery() error = %v", err)
				}
				if values.Get("product_ids") != "1005001234567890,1005001234567891" {
					t.Fatalf("product_ids = %q", values.Get("product_ids"))
				}
				if values.Get("method") != "aliexpress.affiliate.productdetail.get" {
					t.Fatalf("method = %q", values.Get("method"))
				}

				return jsonResponse(`{
			"aliexpress_affiliate_productdetail_get_response": {
				"resp_result": {
					"resp_code": 20010000,
					"resp_msg": "success",
					"result": {
						"current_record_count": 2,
						"products": {
							"product": [
								{
									"product_id": 1005001234567890,
									"product_title": "Keyboard",
									"target_sale_price": "89.99",
									"target_sale_price_currency": "USD",
									"product_main_image_url": "https://img.example/1.jpg",
									"product_detail_url": "https://www.aliexpress.com/item/1005001234567890.html",
									"product_small_image_urls": {"string": ["https://img.example/1-1.jpg"]},
									"promotion_link": "https://s.click.aliexpress.com/promo",
									"target_original_price": "99.99",
									"shop_name": "Gugu Store",
									"sku_id": 20001
								},
								{
									"product_id": 1005001234567891,
									"product_title": "Mouse",
									"target_sale_price": "19.99",
									"target_sale_price_currency": "USD",
									"product_main_image_url": "https://img.example/2.jpg",
									"product_detail_url": "https://www.aliexpress.com/item/1005001234567891.html",
									"shop_name": "Gugu Store"
								}
							]
						}
					}
				}
			}
		}`)
			}),
		},
	})

	result, err := client.GetAffiliateProductDetail(context.Background(), ProductDetailInput{
		ProductIDs: []string{"1005001234567890", "1005001234567891"},
	})
	if err != nil {
		t.Fatalf("GetAffiliateProductDetail() error = %v", err)
	}

	if result.CurrentRecordCount != 2 {
		t.Fatalf("CurrentRecordCount = %d", result.CurrentRecordCount)
	}
	if len(result.Products) != 2 {
		t.Fatalf("len(Products) = %d", len(result.Products))
	}
	if result.Products[0].ProductTitle != "Keyboard" {
		t.Fatalf("ProductTitle = %q", result.Products[0].ProductTitle)
	}
	if result.Products[0].SKUID != 20001 {
		t.Fatalf("SKUID = %d", result.Products[0].SKUID)
	}
}

func TestGetAffiliateProductSKUDetail(t *testing.T) {
	client := newTestClient(t, &Config{
		BaseURL:     "https://api-sg.aliexpress.com",
		AppKey:      "528586",
		AppSecret:   "secret",
		CallbackURL: "https://googoo-client.vercel.app/callback",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != "/sync" {
					t.Fatalf("path = %q, want /sync", r.URL.Path)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll() error = %v", err)
				}
				values, err := url.ParseQuery(string(body))
				if err != nil {
					t.Fatalf("ParseQuery() error = %v", err)
				}
				if values.Get("product_id") != "1005001234567890" {
					t.Fatalf("product_id = %q", values.Get("product_id"))
				}
				if values.Get("sku_ids") != "20001,20002" {
					t.Fatalf("sku_ids = %q", values.Get("sku_ids"))
				}
				if values.Get("method") != "aliexpress.affiliate.product.sku.detail.get" {
					t.Fatalf("method = %q", values.Get("method"))
				}

				return jsonResponse(`{
			"aliexpress_affiliate_product_sku_detail_get_response": {
				"result": {
					"result": {
						"ae_item_info": {
							"product_id": "1005001234567890",
							"title": "Keyboard",
							"original_link": "https://www.aliexpress.com/item/1005001234567890.html",
							"image_link": "https://img.example/1.jpg",
							"store_name": "Gugu Store"
						},
						"ae_item_sku_info": {
							"traffic_sku_info_list": [
								{
									"sku_id": 20001,
									"currency": "USD",
									"sale_price_with_tax": "89.99",
									"sku_image_link": "https://img.example/sku-1.jpg",
									"size": "M"
								}
							]
						},
						"code": 0,
						"success": true
					}
				}
			}
		}`)
			}),
		},
	})

	result, err := client.GetAffiliateProductSKUDetail(context.Background(), ProductSKUDetailInput{
		ProductID:      "1005001234567890",
		ShipToCountry:  "KR",
		TargetCurrency: "USD",
		TargetLanguage: "EN",
		SKUIDs:         []string{"20001", "20002"},
	})
	if err != nil {
		t.Fatalf("GetAffiliateProductSKUDetail() error = %v", err)
	}

	if !result.Success {
		t.Fatal("Success = false, want true")
	}
	if result.ItemInfo.ProductID != "1005001234567890" {
		t.Fatalf("ProductID = %q", result.ItemInfo.ProductID)
	}
	if len(result.SKUInfos) != 1 {
		t.Fatalf("len(SKUInfos) = %d", len(result.SKUInfos))
	}
	if result.SKUInfos[0].SKUID != 20001 {
		t.Fatalf("SKUID = %d", result.SKUInfos[0].SKUID)
	}
}

func TestSignMatchesSDKRule(t *testing.T) {
	client := newTestClient(t, &Config{
		BaseURL:     "https://api-sg.aliexpress.com",
		AppKey:      "528586",
		AppSecret:   "secret",
		CallbackURL: "https://googoo-client.vercel.app/callback",
	})

	values := url.Values{}
	values.Set("app_key", "528586")
	values.Set("sign_method", signMethodSHA256)
	values.Set("timestamp", "1710000000000")
	values.Set("partner_id", defaultPartnerID)
	values.Set("code", "auth-code")

	signature := client.signLegacy("/auth/token/create", values)
	if signature == "" {
		t.Fatal("signature was empty")
	}
	if signature != strings.ToUpper(signature) {
		t.Fatalf("signature = %q, want uppercase", signature)
	}
}

func newTestClient(t *testing.T, cfg *Config) *HTTPClient {
	t.Helper()

	if cfg == nil {
		cfg = &Config{
			BaseURL:     "https://api-sg.aliexpress.com",
			AppKey:      "528586",
			AppSecret:   "secret",
			CallbackURL: "https://googoo-client.vercel.app/callback",
		}
	}

	client, err := NewHTTPClient(*cfg)
	if err != nil {
		t.Fatalf("NewHTTPClient() error = %v", err)
	}
	client.now = func() time.Time {
		return time.UnixMilli(1710000000000)
	}

	return client
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func jsonResponse(body string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}
