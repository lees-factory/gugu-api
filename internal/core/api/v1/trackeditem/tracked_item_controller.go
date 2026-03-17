package trackeditem

import (
	"fmt"
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	trackeditemrequest "github.com/ljj/gugu-api/internal/core/api/v1/trackeditem/request"
	trackeditemresponse "github.com/ljj/gugu-api/internal/core/api/v1/trackeditem/response"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type Controller struct {
	trackedItemService *domaintrackeditem.Service
	productService     *domainproduct.Service
	productCollector   domainproduct.Collector
}

func NewController(
	trackedItemService *domaintrackeditem.Service,
	productService *domainproduct.Service,
	productCollector domainproduct.Collector,
) *Controller {
	return &Controller{
		trackedItemService: trackedItemService,
		productService:     productService,
		productCollector:   productCollector,
	}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Route("/v1/tracked-items", func(r chi.Router) {
		r.Get("/", apiadvice.Wrap(c.List))
		r.Post("/", apiadvice.Wrap(c.Add))
		r.Delete("/{trackedItemID}", apiadvice.Wrap(c.Delete))
	})
}

func (c *Controller) Add(r *stdhttp.Request) (int, any, error) {
	var req trackeditemrequest.AddTrackedItem
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	market := domainproduct.Market(req.ProviderCommerce)
	if !market.IsSupported() {
		return 0, nil, coreerror.New(coreerror.UnsupportedMarket)
	}

	foundProduct, err := c.productService.FindByMarketAndExternalProductID(r.Context(), market, req.ExternalProductID)
	if err != nil {
		return 0, nil, err
	}

	if foundProduct == nil {
		collected, err := c.productCollector.Collect(r.Context(), domainproduct.CollectInput{
			Market:            market,
			ExternalProductID: req.ExternalProductID,
			OriginalURL:       req.OriginalURL,
		})
		if err != nil {
			return 0, nil, fmt.Errorf("collect product: %w", err)
		}

		foundProduct, err = c.productService.Create(r.Context(), domainproduct.CreateInput{
			Market:            collected.Market,
			ExternalProductID: collected.ExternalProductID,
			OriginalURL:       collected.OriginalURL,
			Title:             collected.Title,
			MainImageURL:      collected.MainImageURL,
			CurrentPrice:      collected.CurrentPrice,
			Currency:          collected.Currency,
			ProductURL:        collected.ProductURL,
			CollectionSource:  collected.CollectionSource,
		})
		if err != nil {
			return 0, nil, fmt.Errorf("create product: %w", err)
		}
	}

	result, err := c.trackedItemService.Add(r.Context(), domaintrackeditem.AddInput{
		UserID:      req.UserID,
		ProductID:   foundProduct.ID,
		OriginalURL: req.OriginalURL,
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(
		trackeditemresponse.NewAddTrackedItem(result.TrackedItem, *foundProduct, result.AlreadyTracked),
	), nil
}

func (c *Controller) List(r *stdhttp.Request) (int, any, error) {
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	trackedItems, err := c.trackedItemService.ListByUserID(r.Context(), userID)
	if err != nil {
		return 0, nil, err
	}

	items := make([]trackeditemresponse.ListTrackedItem, 0, len(trackedItems))
	for _, tracked := range trackedItems {
		foundProduct, err := c.productService.FindByID(r.Context(), tracked.ProductID)
		if err != nil {
			return 0, nil, err
		}
		if foundProduct == nil {
			continue
		}
		items = append(items, trackeditemresponse.NewListTrackedItem(tracked, *foundProduct))
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(items), nil
}

func (c *Controller) Delete(r *stdhttp.Request) (int, any, error) {
	trackedItemID := strings.TrimSpace(chi.URLParam(r, "trackedItemID"))
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	if err := c.trackedItemService.Delete(r.Context(), trackedItemID, userID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}
