package integration

import (
	"errors"
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apierror "github.com/ljj/gugu-api/internal/core/api/error"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	alirequest "github.com/ljj/gugu-api/internal/core/api/v1/integration/request"
	aliresponse "github.com/ljj/gugu-api/internal/core/api/v1/integration/response"
	domainintegration "github.com/ljj/gugu-api/internal/core/domain/integration"
)

type AliExpressController struct {
	service       *domainintegration.AliExpressConnectionService
	productClient clientaliexpress.ProductClient
}

func NewAliExpressController(service *domainintegration.AliExpressConnectionService, productClient clientaliexpress.ProductClient) *AliExpressController {
	return &AliExpressController{service: service, productClient: productClient}
}

func (c *AliExpressController) RegisterRoutes(r chi.Router) {
	r.Route("/v1/integrations/aliexpress", func(r chi.Router) {
		r.Post("/authorize-url", apiadvice.Wrap(c.BuildAuthorizationURL))
		r.Post("/exchange-code", apiadvice.Wrap(c.ExchangeCode))
		r.Get("/connection-status", apiadvice.Wrap(c.GetConnectionStatus))
		r.Get("/product-detail", apiadvice.Wrap(c.GetProductDetail))
		r.Get("/product-sku-detail", apiadvice.Wrap(c.GetProductSKUDetail))
	})
}

func (c *AliExpressController) BuildAuthorizationURL(r *stdhttp.Request) (int, any, error) {
	var req alirequest.AliExpressAuthorizationURL
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	result, err := c.service.BuildAuthorizationURL(r.Context(), domainintegration.BuildAliExpressAuthorizationURLInput{
		UserID: req.UserID,
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(aliresponse.NewAliExpressAuthorizationURL(*result)), nil
}

func (c *AliExpressController) ExchangeCode(r *stdhttp.Request) (int, any, error) {
	var req alirequest.AliExpressExchangeCode
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	result, err := c.service.ExchangeCode(r.Context(), domainintegration.ExchangeAliExpressCodeInput{
		UserID: req.UserID,
		Code:   req.Code,
	})
	if err != nil {
		return 0, nil, mapAliExpressError(err)
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(aliresponse.NewAliExpressConnectionStatus(*result)), nil
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

func (c *AliExpressController) GetProductDetail(r *stdhttp.Request) (int, any, error) {
	productID := strings.TrimSpace(r.URL.Query().Get("product_id"))
	if productID == "" {
		return 0, nil, apierror.InvalidRequestError()
	}

	result, err := c.productClient.GetAffiliateProductDetail(r.Context(), clientaliexpress.ProductDetailInput{
		ProductIDs:     []string{productID},
		TargetCurrency: "KRW",
		TargetLanguage: "KO",
		Country:        "KR",
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

	result, err := c.productClient.GetAffiliateProductSKUDetail(r.Context(), clientaliexpress.ProductSKUDetailInput{
		ProductID:      productID,
		ShipToCountry:  "KR",
		TargetCurrency: "KRW",
		TargetLanguage: "KO",
	})
	if err != nil {
		return 0, nil, mapAliExpressError(err)
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(result), nil
}

func (c *AliExpressController) GetConnectionStatus(r *stdhttp.Request) (int, any, error) {
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	result, err := c.service.GetConnectionStatus(r.Context(), userID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(aliresponse.NewAliExpressConnectionStatus(*result)), nil
}
