package pricealert

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

type Controller struct {
	priceAlertService *domainpricealert.Service
}

func NewController(priceAlertService *domainpricealert.Service) *Controller {
	return &Controller{priceAlertService: priceAlertService}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Post("/v1/alerts", apiadvice.Wrap(c.Register))
	r.Get("/v1/alerts", apiadvice.Wrap(c.ListMyAlerts))
	r.Delete("/v1/alerts/{skuID}", apiadvice.Wrap(c.Unregister))
}

func (c *Controller) Register(r *stdhttp.Request) (int, any, error) {
	req, err := request.ParseRegisterPriceAlert(r)
	if err != nil {
		return 0, nil, err
	}

	alert, err := c.priceAlertService.Register(r.Context(), req.User.ID, req.SKUID, req.Channel)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(
		response.NewPriceAlertItem(alert),
	), nil
}

func (c *Controller) Unregister(r *stdhttp.Request) (int, any, error) {
	req := request.ParseUnregisterPriceAlert(r)

	if err := c.priceAlertService.Unregister(r.Context(), req.User.ID, req.SKUID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}

func (c *Controller) ListMyAlerts(r *stdhttp.Request) (int, any, error) {
	req := request.ParseListMyAlerts(r)

	alerts, err := c.priceAlertService.ListByUserID(r.Context(), req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewPriceAlertList(alerts),
	), nil
}
