package product

import (
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	productresponse "github.com/ljj/gugu-api/internal/core/api/v1/product/response"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type Controller struct {
	productService      *domainproduct.Service
	priceHistoryService *domainpricehistory.Service
	trackedItemService  *domaintrackeditem.Service
}

func NewController(
	productService *domainproduct.Service,
	priceHistoryService *domainpricehistory.Service,
	trackedItemService *domaintrackeditem.Service,
) *Controller {
	return &Controller{
		productService:      productService,
		priceHistoryService: priceHistoryService,
		trackedItemService:  trackedItemService,
	}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Route("/v1/products", func(r chi.Router) {
		r.Get("/{productID}", apiadvice.Wrap(c.GetDetail))
	})
}

func (c *Controller) GetDetail(r *stdhttp.Request) (int, any, error) {
	productID := strings.TrimSpace(chi.URLParam(r, "productID"))
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	foundProduct, err := c.productService.FindByID(r.Context(), productID)
	if err != nil {
		return 0, nil, err
	}
	if foundProduct == nil {
		return 0, nil, coreerror.New(coreerror.ProductNotFound)
	}

	histories, err := c.priceHistoryService.ListByProductID(r.Context(), foundProduct.ID)
	if err != nil {
		return 0, nil, err
	}

	isTrackedByUser := false
	trackedItemID := ""
	if userID != "" {
		tracked, err := c.trackedItemService.FindByUserIDAndProductID(r.Context(), userID, foundProduct.ID)
		if err != nil {
			return 0, nil, err
		}
		if tracked != nil {
			isTrackedByUser = true
			trackedItemID = tracked.ID
		}
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(productresponse.NewProductDetail(
		*foundProduct, histories, isTrackedByUser, trackedItemID,
	)), nil
}
