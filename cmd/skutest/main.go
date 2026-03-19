package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	appKey := os.Getenv("ALIEXPRESS_APP_KEY")
	appSecret := os.Getenv("ALIEXPRESS_APP_SECRET")
	if appKey == "" || appSecret == "" {
		log.Fatal("ALIEXPRESS_APP_KEY and ALIEXPRESS_APP_SECRET must be set")
	}

	productID := "1005006851369167"
	if len(os.Args) > 1 {
		productID = os.Args[1]
	}

	baseURL := "https://api-sg.aliexpress.com"

	fmt.Printf("=== SKU Detail for: %s ===\n\n", productID)
	body := callAPI(baseURL, appKey, appSecret, "aliexpress.affiliate.product.sku.detail.get", map[string]string{
		"product_id":      productID,
		"ship_to_country": "KR",
		"target_currency": "KRW",
		"target_language": "KO",
	})

	var raw map[string]any
	json.Unmarshal(body, &raw)

	// Navigate: response > result > result > ae_item_sku_info > traffic_sku_info_list
	skuList := navigateJSON(raw,
		"aliexpress_affiliate_product_sku_detail_get_response",
		"result", "result", "ae_item_sku_info", "traffic_sku_info_list")

	skus, ok := skuList.([]any)
	if !ok {
		fmt.Println("Failed to parse SKU list")
		fmt.Println(string(body))
		return
	}

	fmt.Printf("Total SKUs returned: %d\n\n", len(skus))

	// Group by color
	colorSizes := make(map[string][]string)
	colorIDs := make(map[string][]string)
	for _, s := range skus {
		sku, _ := s.(map[string]any)
		color, _ := sku["color"].(string)
		size, _ := sku["size"].(string)
		skuID := fmt.Sprintf("%.0f", sku["sku_id"])
		colorSizes[color] = append(colorSizes[color], size)
		colorIDs[color] = append(colorIDs[color], skuID)
	}

	fmt.Printf("Unique colors: %d\n", len(colorSizes))
	for color, sizes := range colorSizes {
		ids := colorIDs[color]
		fmt.Printf("  %-12s → sizes: %v\n", color, sizes)
		fmt.Printf("  %-12s   sku_ids: %v\n", "", ids)
	}

	// Print first SKU as sample
	fmt.Println("\n=== Sample SKU (first) ===")
	sample, _ := json.MarshalIndent(skus[0], "", "  ")
	fmt.Println(string(sample))
}

func navigateJSON(data any, keys ...string) any {
	current := data
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = m[key]
	}
	return current
}

func callAPI(baseURL, appKey, appSecret, apiName string, apiParams map[string]string) []byte {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)

	allParams := map[string]string{
		"app_key":     appKey,
		"sign_method": "sha256",
		"timestamp":   ts,
		"method":      apiName,
		"v":           "2.0",
		"format":      "json",
		"partner_id":  "gugu-api",
	}
	for k, v := range apiParams {
		if strings.TrimSpace(v) != "" {
			allParams[k] = v
		}
	}
	allParams["sign"] = signTOP(appSecret, allParams)

	formValues := url.Values{}
	for k, v := range allParams {
		formValues.Set(k, v)
	}

	req, _ := http.NewRequest(http.MethodPost, baseURL+"/sync", strings.NewReader(formValues.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	resp, err := (&http.Client{Timeout: 15 * time.Second}).Do(req)
	if err != nil {
		log.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body
}

func signTOP(appSecret string, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		v := params[k]
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			continue
		}
		b.WriteString(k)
		b.WriteString(v)
	}

	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte(b.String()))
	return strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
}
