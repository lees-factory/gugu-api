package integration

import (
	"context"
	"errors"
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainintegration "github.com/ljj/gugu-api/internal/core/domain/integration"
	apierror "github.com/ljj/gugu-api/internal/core/support/error"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

type AliExpressController struct {
	service       *domainintegration.AliExpressConnectionService
	productClient clientaliexpress.ProductClient
	tokenStore    clientaliexpress.TokenStore
}

func NewAliExpressController(service *domainintegration.AliExpressConnectionService, productClient clientaliexpress.ProductClient, tokenStore clientaliexpress.TokenStore) *AliExpressController {
	return &AliExpressController{service: service, productClient: productClient, tokenStore: tokenStore}
}

func (c *AliExpressController) resolveAccessToken(ctx context.Context) string {
	record, err := c.tokenStore.FindOne(ctx)
	if err != nil || record == nil {
		return ""
	}
	return record.AccessToken
}

func (c *AliExpressController) RegisterRoutes(r chi.Router) {
	r.Route("/v1/integrations/aliexpress", func(r chi.Router) {
		r.Post("/authorize-url", apiadvice.Wrap(c.BuildAuthorizationURL))
		r.Post("/exchange-code", apiadvice.Wrap(c.ExchangeCode))
		r.Post("/refresh-token", apiadvice.Wrap(c.RefreshToken))
		r.Get("/connection-status", apiadvice.Wrap(c.GetConnectionStatus))
		r.Get("/categories", apiadvice.Wrap(c.GetCategories))
		r.Get("/products", apiadvice.Wrap(c.GetProducts))
		r.Get("/product-detail", apiadvice.Wrap(c.GetProductDetail))
		r.Get("/product-sku-detail", apiadvice.Wrap(c.GetProductSKUDetail))
	})
}

func (c *AliExpressController) BuildAuthorizationURL(r *stdhttp.Request) (int, any, error) {
	result, err := c.service.BuildAuthorizationURL(r.Context())
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewAliExpressAuthorizationURL(*result)), nil
}

func (c *AliExpressController) ExchangeCode(r *stdhttp.Request) (int, any, error) {
	var req request.AliExpressExchangeCode
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	result, err := c.service.ExchangeCode(r.Context(), domainintegration.ExchangeAliExpressCodeInput{
		Code: req.Code,
	})
	if err != nil {
		return 0, nil, mapAliExpressError(err)
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewAliExpressConnectionStatus(*result)), nil
}

func mapAliExpressError(err error) error {
	var remoteErr *clientaliexpress.RemoteError
	if errors.As(err, &remoteErr) {
		return &apierror.APIError{
			Status:  stdhttp.StatusBadRequest,
			Code:    "E402",
			Message: remoteErr.Message,
			Data: map[string]string{
				"provider":    "aliexpress",
				"remote_code": remoteErr.Code,
				"request_id":  remoteErr.RequestID,
			},
			Cause: err,
		}
	}

	return err
}

func (c *AliExpressController) GetCategories(r *stdhttp.Request) (int, any, error) {
	result, err := c.productClient.GetAffiliateCategories(r.Context())
	if err != nil {
		return 0, nil, mapAliExpressError(err)
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(result), nil
}

func (c *AliExpressController) GetProducts(r *stdhttp.Request) (int, any, error) {
	q := r.URL.Query()

	result, err := c.productClient.GetAffiliateProducts(r.Context(), clientaliexpress.ProductQueryInput{
		CategoryIDs:    strings.TrimSpace(q.Get("category_ids")),
		Keywords:       strings.TrimSpace(q.Get("keywords")),
		MaxSalePrice:   strings.TrimSpace(q.Get("max_sale_price")),
		MinSalePrice:   strings.TrimSpace(q.Get("min_sale_price")),
		PageNo:         strings.TrimSpace(q.Get("page_no")),
		PageSize:       strings.TrimSpace(q.Get("page_size")),
		Sort:           strings.TrimSpace(q.Get("sort")),
		TargetCurrency: "KRW",
		TargetLanguage: "KO",
		ShipToCountry:  strings.TrimSpace(q.Get("ship_to_country")),
		TrackingID:     strings.TrimSpace(q.Get("tracking_id")),
		AccessToken:    c.resolveAccessToken(r.Context()),
	})
	if err != nil {
		return 0, nil, mapAliExpressError(err)
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(result), nil
}

func (c *AliExpressController) GetProductDetail(r *stdhttp.Request) (int, any, error) {
	productID := strings.TrimSpace(r.URL.Query().Get("product_id"))
	if productID == "" {
		return 0, nil, apierror.InvalidRequestError()
	}

	result, err := c.productClient.GetAffiliateProductDetail(r.Context(), clientaliexpress.ProductDetailInput{
		ProductIDs:     []string{productID},
		TargetCurrency: "KRW",
		TargetLanguage: "KO",
		AccessToken:    c.resolveAccessToken(r.Context()),
	})
	if err != nil {
		return 0, nil, mapAliExpressError(err)
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(result), nil
}

func (c *AliExpressController) GetProductSKUDetail(r *stdhttp.Request) (int, any, error) {
	productID := strings.TrimSpace(r.URL.Query().Get("product_id"))
	if productID == "" {
		return 0, nil, apierror.InvalidRequestError()
	}

	var skuIDs []string
	if raw := strings.TrimSpace(r.URL.Query().Get("sku_ids")); raw != "" {
		skuIDs = strings.Split(raw, ",")
	}

	result, err := c.productClient.GetAffiliateProductSKUDetail(r.Context(), clientaliexpress.ProductSKUDetailInput{
		ProductID:      productID,
		ShipToCountry:  "KR",
		TargetCurrency: "KRW",
		TargetLanguage: "KO",
		SKUIDs:         skuIDs,
		AccessToken:    c.resolveAccessToken(r.Context()),
	})
	if err != nil {
		return 0, nil, mapAliExpressError(err)
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(result), nil
}

func (c *AliExpressController) RefreshToken(r *stdhttp.Request) (int, any, error) {
	result, err := c.service.RefreshToken(r.Context())
	if err != nil {
		return 0, nil, mapAliExpressError(err)
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewAliExpressConnectionStatus(*result)), nil
}

func (c *AliExpressController) GetConnectionStatus(r *stdhttp.Request) (int, any, error) {
	result, err := c.service.GetConnectionStatus(r.Context())
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewAliExpressConnectionStatus(*result)), nil
}
