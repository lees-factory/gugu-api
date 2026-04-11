package product

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

type Controller struct {
	productService     *domainproduct.Service
	snapshotService    *domainps.Service
	trackedItemService *domaintrackeditem.Service
}

func NewController(
	productService *domainproduct.Service,
	snapshotService *domainps.Service,
	trackedItemService *domaintrackeditem.Service,
) *Controller {
	return &Controller{
		productService:     productService,
		snapshotService:    snapshotService,
		trackedItemService: trackedItemService,
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

	currency := "KRW"
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

	skus, err := c.productService.FindSKUsByProductID(r.Context(), foundProduct.ID)
	if err != nil {
		return 0, nil, err
	}
	currentBySKUID := c.resolveSKUCurrentSnapshots(r.Context(), skus, currency)

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewProductDetail(
		*foundProduct, nil, skus, currentBySKUID, isTrackedByUser, trackedItemID, currency,
	)), nil
}

func (c *Controller) GetPriceTrend(r *stdhttp.Request) (int, any, error) {
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
	return stdhttp.StatusBadRequest, nil, coreerror.New(coreerror.SKUIDRequired)
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
	currency := "KRW"
	currentBySKUID := c.resolveSKUCurrentSnapshots(r.Context(), skus, currency)

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewProductSKUsWithCurrentPrice(skus, currentBySKUID),
	), nil
}

func (c *Controller) resolveSKUCurrentSnapshots(ctx context.Context, skus []domainproduct.SKU, currency string) map[string]response.SKUCurrentSnapshot {
	result := make(map[string]response.SKUCurrentSnapshot, len(skus))
	if c.snapshotService == nil {
		return result
	}

	now := time.Now()
	today := now.Format(time.DateOnly)
	from := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, sku := range skus {
		snapshots, err := c.snapshotService.ListSKUSnapshotsByDateRange(ctx, strings.TrimSpace(sku.ID), currency, from, now.Add(24*time.Hour))
		if err != nil {
			continue
		}
		snapshot := resolveLatestSnapshot(snapshots, today)
		if snapshot == nil || strings.TrimSpace(snapshot.Price) == "" {
			continue
		}
		result[strings.TrimSpace(sku.ID)] = response.SKUCurrentSnapshot{
			Price:         strings.TrimSpace(snapshot.Price),
			OriginalPrice: strings.TrimSpace(snapshot.OriginalPrice),
			Currency:      strings.TrimSpace(snapshot.Currency),
		}
	}

	return result
}

func resolveLatestSnapshot(snapshots []domainps.SKUPriceSnapshot, today string) *domainps.SKUPriceSnapshot {
	if len(snapshots) == 0 {
		return nil
	}

	var latestToday *domainps.SKUPriceSnapshot
	var latest *domainps.SKUPriceSnapshot
	for i := range snapshots {
		s := snapshots[i]
		if strings.TrimSpace(s.Price) == "" {
			continue
		}
		if s.SnapshotDate.Format(time.DateOnly) == today {
			if latestToday == nil || s.SnapshotDate.After(latestToday.SnapshotDate) {
				latestToday = &s
			}
		}
		if latest == nil || s.SnapshotDate.After(latest.SnapshotDate) {
			latest = &s
		}
	}

	if latestToday != nil {
		return latestToday
	}
	return latest
}
