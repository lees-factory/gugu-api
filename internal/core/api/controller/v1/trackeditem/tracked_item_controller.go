package trackeditem

import (
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

type Controller struct {
	trackedItemService     *domaintrackeditem.Service
	skuPriceHistoryService *domainsph.Service
	priceAlertService      *domainpricealert.Service
}

func NewController(trackedItemService *domaintrackeditem.Service, skuPriceHistoryService *domainsph.Service, priceAlertService *domainpricealert.Service) *Controller {
	return &Controller{
		trackedItemService:     trackedItemService,
		skuPriceHistoryService: skuPriceHistoryService,
		priceAlertService:      priceAlertService,
	}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Route("/v1/tracked-items", func(r chi.Router) {
		r.Get("/", apiadvice.Wrap(c.List))
		r.Post("/", apiadvice.Wrap(c.Add))
		r.Get("/{trackedItemID}", apiadvice.Wrap(c.GetDetail))
		r.Get("/{trackedItemID}/price-alert", apiadvice.Wrap(c.GetPriceAlert))
		r.Delete("/{trackedItemID}", apiadvice.Wrap(c.Delete))
		r.Patch("/{trackedItemID}/sku", apiadvice.Wrap(c.SelectSKU))
		r.Get("/{trackedItemID}/sku-price-histories", apiadvice.Wrap(c.GetSKUPriceHistories))
	})
}

func (c *Controller) Add(r *stdhttp.Request) (int, any, error) {
	req, err := request.ParseAddTrackedItems(r)
	if err != nil {
		return stdhttp.StatusBadRequest, nil, err
	}

	results := make([]domaintrackeditem.AddTrackedItemResult, 0, len(req.Items))
	for _, item := range req.Items {
		result, err := c.trackedItemService.AddTrackedItem(r.Context(), domaintrackeditem.AddTrackedItemInput{
			UserID:            req.User.ID,
			ProviderCommerce:  item.ProviderCommerce,
			ExternalProductID: item.ExternalProductID,
			OriginalURL:       item.OriginalURL,
			Currency:          item.Currency,
			Language:          item.Language,
		})
		if err != nil {
			return 0, nil, err
		}
		results = append(results, *result)
	}

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(
		response.NewAddTrackedItems(results),
	), nil
}

func (c *Controller) GetDetail(r *stdhttp.Request) (int, any, error) {
	req := request.ParseGetTrackedItemDetail(r)

	detail, err := c.trackedItemService.GetDetail(r.Context(), req.TrackedItemID, req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewTrackedItemDetail(detail),
	), nil
}

func (c *Controller) List(r *stdhttp.Request) (int, any, error) {
	req := request.ParseListTrackedItems(r)

	result, err := c.trackedItemService.ListWithProductsCursor(r.Context(), req.User.ID, req.Cursor)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewListTrackedItemsPage(result),
	), nil
}

func (c *Controller) Delete(r *stdhttp.Request) (int, any, error) {
	req := request.ParseDeleteTrackedItem(r)

	if err := c.trackedItemService.Delete(r.Context(), req.TrackedItemID, req.User.ID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}

func (c *Controller) GetSKUPriceHistories(r *stdhttp.Request) (int, any, error) {
	req := request.ParseGetSKUPriceHistories(r)

	if req.SKUID == "" {
		return stdhttp.StatusBadRequest, nil, coreerror.New(coreerror.SKUIDRequired)
	}

	found, err := c.trackedItemService.GetDetail(r.Context(), req.TrackedItemID, req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	currency := resolveSKUPriceHistoryCurrency(req.Currency, found.TrackedItem.Currency)

	histories, err := c.skuPriceHistoryService.ListBySKUID(r.Context(), req.SKUID, currency)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewSKUPriceHistories(histories),
	), nil
}

func (c *Controller) GetPriceAlert(r *stdhttp.Request) (int, any, error) {
	req := request.ParseGetTrackedItemPriceAlert(r)

	detail, err := c.trackedItemService.GetDetail(r.Context(), req.TrackedItemID, req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	skuID := strings.TrimSpace(req.SKUID)
	if skuID == "" {
		skuID = strings.TrimSpace(detail.TrackedItem.SKUID)
	}
	if skuID == "" {
		return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewPriceAlertState(nil)), nil
	}

	if !containsSKUID(detail.SKUs, skuID) {
		return stdhttp.StatusBadRequest, nil, coreerror.New(coreerror.SKUIDRequired)
	}

	alertState := c.resolvePriceAlertStateBySKUID(r, detail.TrackedItem.UserID, skuID)
	return stdhttp.StatusOK, apiresponse.SuccessWithData(*alertState), nil
}

func resolveSKUPriceHistoryCurrency(requested string, trackedItemCurrency string) string {
	currency := strings.ToUpper(strings.TrimSpace(requested))
	if currency != "" {
		return currency
	}

	currency = strings.ToUpper(strings.TrimSpace(trackedItemCurrency))
	if currency != "" {
		return currency
	}

	return "KRW"
}

func (c *Controller) resolvePriceAlertStateBySKUID(r *stdhttp.Request, userID string, skuID string) *response.PriceAlertState {
	if c.priceAlertService == nil || strings.TrimSpace(skuID) == "" {
		defaultState := response.NewPriceAlertState(nil)
		return &defaultState
	}

	alert, err := c.priceAlertService.FindByUserIDAndSKUID(r.Context(), userID, skuID)
	if err != nil {
		defaultState := response.NewPriceAlertState(nil)
		return &defaultState
	}
	state := response.NewPriceAlertState(alert)
	return &state
}

func containsSKUID(skus []domainproduct.SKU, skuID string) bool {
	for _, sku := range skus {
		if sku.ID == skuID {
			return true
		}
	}
	return false
}

func (c *Controller) SelectSKU(r *stdhttp.Request) (int, any, error) {
	req, err := request.ParseSelectSKU(r)
	if err != nil {
		return 0, nil, err
	}

	if err := c.trackedItemService.SelectSKU(r.Context(), req.TrackedItemID, req.User.ID, req.SKUID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}
