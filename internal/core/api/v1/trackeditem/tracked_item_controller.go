package trackeditem

import (
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	trackeditemrequest "github.com/ljj/gugu-api/internal/core/api/v1/trackeditem/request"
	trackeditemresponse "github.com/ljj/gugu-api/internal/core/api/v1/trackeditem/response"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type Controller struct {
	trackedItemService *domaintrackeditem.Service
}

func NewController(trackedItemService *domaintrackeditem.Service) *Controller {
	return &Controller{trackedItemService: trackedItemService}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Route("/v1/tracked-items", func(r chi.Router) {
		r.Get("/", apiadvice.Wrap(c.List))
		r.Post("/", apiadvice.Wrap(c.Add))
		r.Get("/{trackedItemID}", apiadvice.Wrap(c.GetDetail))
		r.Delete("/{trackedItemID}", apiadvice.Wrap(c.Delete))
		r.Patch("/{trackedItemID}/sku", apiadvice.Wrap(c.SelectSKU))
	})
}

func (c *Controller) Add(r *stdhttp.Request) (int, any, error) {
	var req trackeditemrequest.AddTrackedItem
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	result, err := c.trackedItemService.AddTrackedItem(r.Context(), domaintrackeditem.AddTrackedItemInput{
		UserID:            req.UserID,
		ProviderCommerce:  req.ProviderCommerce,
		ExternalProductID: req.ExternalProductID,
		OriginalURL:       req.OriginalURL,
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(
		trackeditemresponse.NewAddTrackedItemFromResult(result),
	), nil
}

func (c *Controller) GetDetail(r *stdhttp.Request) (int, any, error) {
	trackedItemID := strings.TrimSpace(chi.URLParam(r, "trackedItemID"))
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	detail, err := c.trackedItemService.GetDetail(r.Context(), trackedItemID, userID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		trackeditemresponse.NewTrackedItemDetail(detail),
	), nil
}

func (c *Controller) List(r *stdhttp.Request) (int, any, error) {
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	items, err := c.trackedItemService.ListWithProducts(r.Context(), userID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		trackeditemresponse.NewListTrackedItems(items),
	), nil
}

func (c *Controller) Delete(r *stdhttp.Request) (int, any, error) {
	trackedItemID := strings.TrimSpace(chi.URLParam(r, "trackedItemID"))
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	if err := c.trackedItemService.Delete(r.Context(), trackedItemID, userID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}

func (c *Controller) SelectSKU(r *stdhttp.Request) (int, any, error) {
	trackedItemID := strings.TrimSpace(chi.URLParam(r, "trackedItemID"))

	var req trackeditemrequest.SelectSKU
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	if err := c.trackedItemService.SelectSKU(r.Context(), trackedItemID, req.UserID, req.SKUID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}
