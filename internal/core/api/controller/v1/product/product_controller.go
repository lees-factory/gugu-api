package product

import (
	"fmt"
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

type Controller struct {
	productService      *domainproduct.Service
	priceHistoryService *domainpricehistory.Service
	snapshotService     *domainps.Service
	trackedItemService  *domaintrackeditem.Service
}

func NewController(
	productService *domainproduct.Service,
	priceHistoryService *domainpricehistory.Service,
	snapshotService *domainps.Service,
	trackedItemService *domaintrackeditem.Service,
) *Controller {
	return &Controller{
		productService:      productService,
		priceHistoryService: priceHistoryService,
		snapshotService:     snapshotService,
		trackedItemService:  trackedItemService,
	}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Route("/v1/products", func(r chi.Router) {
		r.Get("/{productID}", apiadvice.Wrap(c.GetDetail))
		r.Get("/{productID}/skus", apiadvice.Wrap(c.ListSKUs))
		r.Get("/{productID}/price-trend", apiadvice.Wrap(c.GetPriceTrend))
	})
}

func (c *Controller) GetDetail(r *stdhttp.Request) (int, any, error) {
	req := request.ParseGetProductDetail(r)

	foundProduct, err := c.productService.FindByID(r.Context(), req.ProductID)
	if err != nil {
		return 0, nil, err
	}

	currency := foundProduct.Currency
	isTrackedByUser := false
	trackedItemID := ""
	if req.User.ID != "" {
		tracked, err := c.trackedItemService.FindByUserIDAndProductID(r.Context(), req.User.ID, foundProduct.ID)
		if err != nil {
			return 0, nil, err
		}
		if tracked != nil {
			isTrackedByUser = true
			trackedItemID = tracked.ID
			if tracked.Currency != "" {
				currency = tracked.Currency
			}
		}
	}
	if currency == "" {
		currency = "KRW"
	}

	histories, err := c.priceHistoryService.ListByProductID(r.Context(), foundProduct.ID, currency)
	if err != nil {
		return 0, nil, err
	}

	skus, err := c.productService.FindSKUsByProductID(r.Context(), foundProduct.ID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewProductDetail(
		*foundProduct, histories, skus, isTrackedByUser, trackedItemID, currency,
	)), nil
}

func (c *Controller) GetPriceTrend(r *stdhttp.Request) (int, any, error) {
	productID := chi.URLParam(r, "productID")
	skuID := r.URL.Query().Get("sku_id")
	currency := r.URL.Query().Get("currency")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if currency == "" {
		currency = "KRW"
	}

	from, err := time.Parse(time.DateOnly, fromStr)
	if err != nil {
		return stdhttp.StatusBadRequest, nil, fmt.Errorf("invalid 'from' date format, expected YYYY-MM-DD: %w", err)
	}
	to, err := time.Parse(time.DateOnly, toStr)
	if err != nil {
		return stdhttp.StatusBadRequest, nil, fmt.Errorf("invalid 'to' date format, expected YYYY-MM-DD: %w", err)
	}

	if skuID != "" {
		snapshots, err := c.snapshotService.ListSKUSnapshotsByDateRange(r.Context(), skuID, currency, from, to)
		if err != nil {
			return 0, nil, err
		}
		return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewSKUPriceTrend(snapshots)), nil
	}

	snapshots, err := c.snapshotService.ListProductSnapshotsByDateRange(r.Context(), productID, currency, from, to)
	if err != nil {
		return 0, nil, err
	}
	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewProductPriceTrend(snapshots)), nil
}

func (c *Controller) ListSKUs(r *stdhttp.Request) (int, any, error) {
	req := request.ParseListProductSKUs(r)

	foundProduct, err := c.productService.FindByID(r.Context(), req.ProductID)
	if err != nil {
		return 0, nil, err
	}

	skus, err := c.productService.FindSKUsByProductID(r.Context(), foundProduct.ID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewProductSKUs(skus),
	), nil
}
