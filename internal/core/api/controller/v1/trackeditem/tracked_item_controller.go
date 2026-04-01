package trackeditem

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

type Controller struct {
	trackedItemService     *domaintrackeditem.Service
	skuPriceHistoryService *domainsph.Service
}

func NewController(trackedItemService *domaintrackeditem.Service, skuPriceHistoryService *domainsph.Service) *Controller {
	return &Controller{
		trackedItemService:     trackedItemService,
		skuPriceHistoryService: skuPriceHistoryService,
	}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Route("/v1/tracked-items", func(r chi.Router) {
		r.Get("/", apiadvice.Wrap(c.List))
		r.Post("/", apiadvice.Wrap(c.Add))
		r.Get("/{trackedItemID}", apiadvice.Wrap(c.GetDetail))
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

	currency := req.Currency
	if currency == "" {
		currency = found.TrackedItem.Currency
	}
	if currency == "" {
		currency = "KRW"
	}

	histories, err := c.skuPriceHistoryService.ListBySKUID(r.Context(), req.SKUID, currency)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewSKUPriceHistories(histories),
	), nil
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
