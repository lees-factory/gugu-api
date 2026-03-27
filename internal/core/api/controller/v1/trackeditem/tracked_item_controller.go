package trackeditem

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
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
	req, err := request.ParseAddTrackedItem(r)
	if err != nil {
		return 0, nil, err
	}

	result, err := c.trackedItemService.AddTrackedItem(r.Context(), domaintrackeditem.AddTrackedItemInput{
		UserID:            req.User.ID,
		ProviderCommerce:  req.ProviderCommerce,
		ExternalProductID: req.ExternalProductID,
		OriginalURL:       req.OriginalURL,
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(
		response.NewAddTrackedItemFromResult(result),
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

	items, err := c.trackedItemService.ListWithProducts(r.Context(), req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewListTrackedItems(items),
	), nil
}

func (c *Controller) Delete(r *stdhttp.Request) (int, any, error) {
	req := request.ParseDeleteTrackedItem(r)

	if err := c.trackedItemService.Delete(r.Context(), req.TrackedItemID, req.User.ID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
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
