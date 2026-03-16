package trackeditem

import (
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	trackeditemrequest "github.com/ljj/gugu-api/internal/core/api/v1/trackeditem/request"
	trackeditemresponse "github.com/ljj/gugu-api/internal/core/api/v1/trackeditem/response"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	trackeditemlist "github.com/ljj/gugu-api/internal/core/domain/trackeditemlist"
)

type Controller struct {
	service     *domaintrackeditem.Service
	listService *trackeditemlist.Service
}

func NewController(service *domaintrackeditem.Service, listService *trackeditemlist.Service) *Controller {
	return &Controller{
		service:     service,
		listService: listService,
	}
}

func (c *Controller) Add(r *stdhttp.Request) (int, any, error) {
	var req trackeditemrequest.AddTrackedItem
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	result, err := c.service.Add(r.Context(), domaintrackeditem.AddInput{
		UserID:            req.UserID,
		OriginalURL:       req.OriginalURL,
		Market:            domainproduct.Market(req.ProviderCommerce),
		ExternalProductID: req.ExternalProductID,
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(trackeditemresponse.NewAddTrackedItem(*result)), nil
}

func (c *Controller) List(r *stdhttp.Request) (int, any, error) {
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	items, err := c.listService.List(r.Context(), userID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(trackeditemresponse.NewListTrackedItems(items)), nil
}

func (c *Controller) Delete(r *stdhttp.Request) (int, any, error) {
	trackedItemID := strings.TrimSpace(chi.URLParam(r, "trackedItemID"))
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	if err := c.service.Delete(r.Context(), trackedItemID, userID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}
