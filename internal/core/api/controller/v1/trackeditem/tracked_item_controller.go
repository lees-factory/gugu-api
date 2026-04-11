package trackeditem

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
	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

type Controller struct {
	trackedItemService     *domaintrackeditem.Service
	skuPriceHistoryService *domainsph.Service
	snapshotService        *domainps.Service
	priceAlertService      *domainpricealert.Service
}

func NewController(
	trackedItemService *domaintrackeditem.Service,
	skuPriceHistoryService *domainsph.Service,
	snapshotService *domainps.Service,
	priceAlertService *domainpricealert.Service,
) *Controller {
	return &Controller{
		trackedItemService:     trackedItemService,
		skuPriceHistoryService: skuPriceHistoryService,
		snapshotService:        snapshotService,
		priceAlertService:      priceAlertService,
	}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Route("/v1/tracked-items", func(r chi.Router) {
		r.Get("/", apiadvice.Wrap(c.List))
		r.Post("/", apiadvice.Wrap(c.Add))
		r.Get("/{trackedItemID}", apiadvice.Wrap(c.GetDetail))
		r.Get("/{trackedItemID}/price-alert", apiadvice.Wrap(c.GetPriceAlert))
		r.Post("/{trackedItemID}/price-alert", apiadvice.Wrap(c.RegisterPriceAlert))
		r.Delete("/{trackedItemID}/price-alert", apiadvice.Wrap(c.UnregisterPriceAlert))
		r.Get("/{trackedItemID}/skus/{skuID}/price-alert", apiadvice.Wrap(c.GetPriceAlert))
		r.Post("/{trackedItemID}/skus/{skuID}/price-alert", apiadvice.Wrap(c.RegisterPriceAlert))
		r.Delete("/{trackedItemID}/skus/{skuID}/price-alert", apiadvice.Wrap(c.UnregisterPriceAlert))
		r.Delete("/{trackedItemID}", apiadvice.Wrap(c.Delete))
		r.Patch("/{trackedItemID}/sku", apiadvice.Wrap(c.SelectSKU))
		r.Patch("/{trackedItemID}/language", apiadvice.Wrap(c.UpdateLanguage))
		r.Get("/{trackedItemID}/sku-price-histories", apiadvice.Wrap(c.GetSKUPriceHistories))
		r.Get("/{trackedItemID}/sku-price-trend", apiadvice.Wrap(c.GetSKUPriceTrend))
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
			OriginProductID:   item.OriginProductID,
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
	skuCurrentSnapshots := c.resolveTrackedItemSKUCurrentSnapshots(r.Context(), detail)
	currentPrice := resolveTrackedItemCurrentPrice(detail, skuCurrentSnapshots)
	skus := response.NewProductSKUsWithCurrentPrice(detail.SKUs, skuCurrentSnapshots)

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewTrackedItemDetail(detail, currentPrice, skus),
	), nil
}

func (c *Controller) List(r *stdhttp.Request) (int, any, error) {
	req := request.ParseListTrackedItems(r)

	result, err := c.trackedItemService.ListWithProductsCursor(r.Context(), req.User.ID, req.Cursor)
	if err != nil {
		return 0, nil, err
	}
	currentPriceByTrackedItemID := c.resolveTrackedItemsCurrentPrices(r.Context(), result.Items)

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewListTrackedItemsPage(result, currentPriceByTrackedItemID),
	), nil
}

func (c *Controller) Delete(r *stdhttp.Request) (int, any, error) {
	req := request.ParseDeleteTrackedItem(r)

	if err := c.trackedItemService.Delete(r.Context(), req.TrackedItemID, req.User.ID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}

func (c *Controller) UpdateLanguage(r *stdhttp.Request) (int, any, error) {
	req, err := request.ParseUpdateTrackedItemLanguage(r)
	if err != nil {
		return stdhttp.StatusBadRequest, nil, err
	}

	if err := c.trackedItemService.UpdatePreferredLanguage(r.Context(), req.TrackedItemID, req.User.ID, req.Language); err != nil {
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

func (c *Controller) GetSKUPriceTrend(r *stdhttp.Request) (int, any, error) {
	req := request.ParseGetSKUPriceTrend(r)

	from, err := time.Parse(time.DateOnly, req.From)
	if err != nil {
		return stdhttp.StatusBadRequest, nil, fmt.Errorf("invalid 'from' date format, expected YYYY-MM-DD: %w", err)
	}
	to, err := time.Parse(time.DateOnly, req.To)
	if err != nil {
		return stdhttp.StatusBadRequest, nil, fmt.Errorf("invalid 'to' date format, expected YYYY-MM-DD: %w", err)
	}

	found, err := c.trackedItemService.GetDetail(r.Context(), req.TrackedItemID, req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	skuID, err := resolveTrackedItemPriceAlertSKUID(found, req.SKUID, true)
	if err != nil {
		return 0, nil, err
	}

	currency := resolveSKUPriceHistoryCurrency(req.Currency, found.TrackedItem.Currency)

	snapshots, err := c.snapshotService.ListSKUSnapshotsByDateRange(r.Context(), skuID, currency, from, to)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewSKUPriceTrend(snapshots),
	), nil
}

func (c *Controller) GetPriceAlert(r *stdhttp.Request) (int, any, error) {
	req := request.ParseGetTrackedItemPriceAlert(r)

	detail, err := c.trackedItemService.GetDetail(r.Context(), req.TrackedItemID, req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	skuID, err := resolveTrackedItemPriceAlertSKUID(detail, req.SKUID, false)
	if err != nil {
		return 0, nil, err
	}
	if skuID == "" {
		alertState := c.resolvePriceAlertStateByTrackedItem(r, detail)
		return stdhttp.StatusOK, apiresponse.SuccessWithData(*alertState), nil
	}

	alertState := c.resolvePriceAlertStateBySKUID(r, detail.TrackedItem.UserID, skuID)
	return stdhttp.StatusOK, apiresponse.SuccessWithData(*alertState), nil
}

func (c *Controller) RegisterPriceAlert(r *stdhttp.Request) (int, any, error) {
	req, err := request.ParseRegisterTrackedItemPriceAlert(r)
	if err != nil {
		return stdhttp.StatusBadRequest, nil, err
	}

	detail, err := c.trackedItemService.GetDetail(r.Context(), req.TrackedItemID, req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	skuID, err := resolveTrackedItemPriceAlertSKUID(detail, req.SKUID, true)
	if err != nil {
		return 0, nil, err
	}

	alert, err := c.priceAlertService.Register(r.Context(), req.User.ID, skuID, req.Channel)
	if err != nil {
		return 0, nil, err
	}
	// 등록된 SKU를 tracked item 기본 선택값으로 동기화해 조회 일관성을 맞춘다.
	_ = c.trackedItemService.SelectSKU(r.Context(), req.TrackedItemID, req.User.ID, skuID)

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(response.NewPriceAlertState(alert)), nil
}

func (c *Controller) UnregisterPriceAlert(r *stdhttp.Request) (int, any, error) {
	req := request.ParseUnregisterTrackedItemPriceAlert(r)

	detail, err := c.trackedItemService.GetDetail(r.Context(), req.TrackedItemID, req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	skuID, err := resolveTrackedItemPriceAlertSKUID(detail, req.SKUID, true)
	if err != nil {
		return 0, nil, err
	}

	if err := c.priceAlertService.Unregister(r.Context(), req.User.ID, skuID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
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

func resolveTrackedItemPriceAlertSKUID(detail *domaintrackeditem.TrackedItemDetail, requestedSKUID string, required bool) (string, error) {
	skuID := strings.TrimSpace(requestedSKUID)
	if skuID == "" {
		skuID = strings.TrimSpace(detail.TrackedItem.SKUID)
	}
	if skuID == "" {
		if required {
			return "", coreerror.New(coreerror.SKUIDRequired)
		}
		return "", nil
	}
	if !containsSKUID(detail.SKUs, skuID) {
		return "", coreerror.New(coreerror.SKUIDRequired)
	}
	return skuID, nil
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

func (c *Controller) resolvePriceAlertStateByTrackedItem(r *stdhttp.Request, detail *domaintrackeditem.TrackedItemDetail) *response.PriceAlertState {
	if c.priceAlertService == nil || detail == nil {
		defaultState := response.NewPriceAlertState(nil)
		return &defaultState
	}

	seen := make(map[string]struct{}, len(detail.SKUs)+1)
	candidates := make([]string, 0, len(detail.SKUs)+1)

	selected := strings.TrimSpace(detail.TrackedItem.SKUID)
	if selected != "" {
		seen[selected] = struct{}{}
		candidates = append(candidates, selected)
	}
	for _, sku := range detail.SKUs {
		skuID := strings.TrimSpace(sku.ID)
		if skuID == "" {
			continue
		}
		if _, ok := seen[skuID]; ok {
			continue
		}
		seen[skuID] = struct{}{}
		candidates = append(candidates, skuID)
	}

	for _, skuID := range candidates {
		alert, err := c.priceAlertService.FindByUserIDAndSKUID(r.Context(), detail.TrackedItem.UserID, skuID)
		if err != nil || alert == nil {
			continue
		}
		state := response.NewPriceAlertState(alert)
		return &state
	}

	defaultState := response.NewPriceAlertState(nil)
	return &defaultState
}

func containsSKUID(skus []domainproduct.SKU, skuID string) bool {
	for _, sku := range skus {
		if sku.ID == skuID {
			return true
		}
	}
	return false
}

func (c *Controller) resolveTrackedItemSKUCurrentSnapshots(ctx context.Context, detail *domaintrackeditem.TrackedItemDetail) map[string]response.SKUCurrentSnapshot {
	result := make(map[string]response.SKUCurrentSnapshot, len(detail.SKUs))
	if c.snapshotService == nil {
		return result
	}

	currency := resolveSKUPriceHistoryCurrency("", detail.TrackedItem.Currency)
	now := time.Now()
	today := now.Format(time.DateOnly)

	for _, sku := range detail.SKUs {
		snapshots, err := c.snapshotService.ListSKUSnapshotsByDateRange(
			ctx,
			strings.TrimSpace(sku.ID),
			currency,
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			now.Add(24*time.Hour),
		)
		if err != nil {
			continue
		}
		price := resolveSnapshotCurrentPrice(snapshots, today)
		snapshot := resolveLatestSnapshot(snapshots, today)
		if snapshot == nil || strings.TrimSpace(price) == "" {
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

func (c *Controller) resolveTrackedItemsCurrentPrices(ctx context.Context, items []domaintrackeditem.TrackedItemWithProduct) map[string]string {
	result := make(map[string]string, len(items))
	if c.snapshotService == nil {
		return result
	}

	now := time.Now()
	today := now.Format(time.DateOnly)
	from := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, item := range items {
		skuID := strings.TrimSpace(item.TrackedItem.SKUID)
		if skuID == "" {
			continue
		}
		currency := resolveSKUPriceHistoryCurrency("", item.TrackedItem.Currency)
		snapshots, err := c.snapshotService.ListSKUSnapshotsByDateRange(ctx, skuID, currency, from, now.Add(24*time.Hour))
		if err != nil {
			continue
		}
		price := resolveSnapshotCurrentPrice(snapshots, today)
		if price == "" {
			continue
		}
		result[strings.TrimSpace(item.TrackedItem.ID)] = price
	}

	return result
}

func resolveTrackedItemCurrentPrice(detail *domaintrackeditem.TrackedItemDetail, skuCurrentSnapshots map[string]response.SKUCurrentSnapshot) string {
	skuID, err := resolveTrackedItemPriceAlertSKUID(detail, "", false)
	if err != nil || skuID == "" {
		return ""
	}
	return strings.TrimSpace(skuCurrentSnapshots[strings.TrimSpace(skuID)].Price)
}

func resolveSnapshotCurrentPrice(snapshots []domainps.SKUPriceSnapshot, today string) string {
	snapshot := resolveLatestSnapshot(snapshots, today)
	if snapshot == nil {
		return ""
	}
	return strings.TrimSpace(snapshot.Price)
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
	if latest == nil {
		return nil
	}
	return latest
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
