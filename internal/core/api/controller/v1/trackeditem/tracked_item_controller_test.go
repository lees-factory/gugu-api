package trackeditem

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
	memorypricealert "github.com/ljj/gugu-api/internal/storage/memory/pricealert"
	memorypricesnapshot "github.com/ljj/gugu-api/internal/storage/memory/pricesnapshot"
	memoryproduct "github.com/ljj/gugu-api/internal/storage/memory/product"
	memoryskupricehistory "github.com/ljj/gugu-api/internal/storage/memory/skupricehistory"
	memorytrackeditem "github.com/ljj/gugu-api/internal/storage/memory/trackeditem"
)

func TestResolveSKUPriceHistoryCurrency_UsesOverride(t *testing.T) {
	got := resolveSKUPriceHistoryCurrency("usd", "KRW")
	if got != "USD" {
		t.Fatalf("resolveSKUPriceHistoryCurrency() = %q, want USD", got)
	}
}

func TestResolveSKUPriceHistoryCurrency_UsesTrackedItemCurrencyByDefault(t *testing.T) {
	got := resolveSKUPriceHistoryCurrency("", "krw")
	if got != "KRW" {
		t.Fatalf("resolveSKUPriceHistoryCurrency() = %q, want KRW", got)
	}
}

func TestResolveSKUPriceHistoryCurrency_FallsBackToKRW(t *testing.T) {
	got := resolveSKUPriceHistoryCurrency("", "")
	if got != "KRW" {
		t.Fatalf("resolveSKUPriceHistoryCurrency() = %q, want KRW", got)
	}
}

func TestResolvePriceAlertStateBySKUID_DefaultsOffWithoutService(t *testing.T) {
	controller := &Controller{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tracked-items/tracked-1/price-alert", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}

	got := controller.resolvePriceAlertStateBySKUID(req, "user-1", "sku-1")

	if got == nil {
		t.Fatal("resolvePriceAlertStateBySKUID() returned nil")
	}
	if *got != (response.PriceAlertState{Enabled: false}) {
		t.Fatalf("resolvePriceAlertStateBySKUID() = %+v, want disabled default state", *got)
	}
}

func TestResolvePriceAlertStateBySKUID_ReturnsStoredAlert(t *testing.T) {
	repo := memorypricealert.NewRepository()
	service := domainpricealert.NewService(
		domainpricealert.NewFinder(repo),
		repo,
		testAlertIDGenerator{},
		testAlertClock{},
	)
	if _, err := service.Register(context.Background(), "user-1", "sku-1", "EMAIL"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	controller := &Controller{priceAlertService: service}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tracked-items/tracked-1/price-alert?sku_id=sku-1", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}

	got := controller.resolvePriceAlertStateBySKUID(req, "user-1", "sku-1")

	if got == nil {
		t.Fatal("resolvePriceAlertStateBySKUID() returned nil")
	}
	if !got.Enabled {
		t.Fatalf("resolvePriceAlertStateBySKUID().Enabled = false, want true")
	}
	if got.Channel != "EMAIL" {
		t.Fatalf("resolvePriceAlertStateBySKUID().Channel = %q, want EMAIL", got.Channel)
	}
}

func TestResolveTrackedItemPriceAlertSKUID_UsesRequestedSKU(t *testing.T) {
	detail := &domaintrackeditem.TrackedItemDetail{
		TrackedItem: domaintrackeditem.TrackedItem{SKUID: "sku-selected"},
		SKUs: []domainproduct.SKU{
			{ID: "sku-selected"},
			{ID: "sku-requested"},
		},
	}

	got, err := resolveTrackedItemPriceAlertSKUID(detail, "sku-requested", true)
	if err != nil {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() error = %v", err)
	}
	if got != "sku-requested" {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() = %q, want sku-requested", got)
	}
}

func TestResolveTrackedItemPriceAlertSKUID_UsesSelectedSKUByDefault(t *testing.T) {
	detail := &domaintrackeditem.TrackedItemDetail{
		TrackedItem: domaintrackeditem.TrackedItem{SKUID: "sku-selected"},
		SKUs:        []domainproduct.SKU{{ID: "sku-selected"}},
	}

	got, err := resolveTrackedItemPriceAlertSKUID(detail, "", true)
	if err != nil {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() error = %v", err)
	}
	if got != "sku-selected" {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() = %q, want sku-selected", got)
	}
}

func TestResolveTrackedItemPriceAlertSKUID_AllowsMissingSKUForRead(t *testing.T) {
	detail := &domaintrackeditem.TrackedItemDetail{}

	got, err := resolveTrackedItemPriceAlertSKUID(detail, "", false)
	if err != nil {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() error = %v", err)
	}
	if got != "" {
		t.Fatalf("resolveTrackedItemPriceAlertSKUID() = %q, want empty", got)
	}
}

func TestController_RegisterRoutes_IncludesSKUPriceTrend(t *testing.T) {
	controller := &Controller{}

	chiRouter := chi.NewRouter()
	controller.RegisterRoutes(chiRouter)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/tracked-items/tracked-1/sku-price-trend?sku_id=sku-1&from=2026-04-01&to=2026-04-06", nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}

	routeContext := chi.NewRouteContext()
	if !chiRouter.Match(routeContext, http.MethodGet, req.URL.Path) {
		t.Fatalf("RegisterRoutes() did not register sku-price-trend route")
	}
}

func TestNewController_AssignsSnapshotService(t *testing.T) {
	snapshotService := domainps.NewService(nil)

	controller := NewController(nil, nil, snapshotService, nil)

	if controller.snapshotService != snapshotService {
		t.Fatalf("NewController() did not assign snapshotService")
	}
}

func TestGetSKUPriceHistories_UsesTrackedItemCurrencyByDefault(t *testing.T) {
	controller, historyRepo, _ := newTestTrackedItemController(t, "USD")

	if err := historyRepo.Create(context.Background(), domainsph.SKUPriceHistory{
		SKUID:       "sku-1",
		Price:       "1000",
		Currency:    "KRW",
		RecordedAt:  time.Date(2026, 4, 6, 10, 0, 0, 0, time.UTC),
		ChangeValue: "0",
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := newTrackedItemRequest(t, http.MethodGet, "/v1/tracked-items/tracked-1/sku-price-histories?sku_id=sku-1", "tracked-1")

	status, body, err := controller.GetSKUPriceHistories(req)
	if err != nil {
		t.Fatalf("GetSKUPriceHistories() error = %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("GetSKUPriceHistories() status = %d, want 200", status)
	}

	resp, ok := body.(apiresponse.APIResponse[[]response.SKUPriceHistoryItem])
	if !ok {
		t.Fatalf("GetSKUPriceHistories() body type = %T", body)
	}
	if resp.Data == nil {
		t.Fatal("GetSKUPriceHistories() returned nil data")
	}
	if len(*resp.Data) != 0 {
		t.Fatalf("GetSKUPriceHistories() items = %d, want 0 when tracked item currency defaults to USD", len(*resp.Data))
	}
}

func TestGetSKUPriceTrend_UsesTrackedItemCurrencyByDefault(t *testing.T) {
	controller, _, snapshotRepo := newTestTrackedItemController(t, "KRW")

	if err := snapshotRepo.Upsert(context.Background(), domainps.SKUPriceSnapshot{
		SKUID:         "sku-1",
		SnapshotDate:  time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC),
		Price:         "1100",
		OriginalPrice: "1200",
		Currency:      "KRW",
	}); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if err := snapshotRepo.Upsert(context.Background(), domainps.SKUPriceSnapshot{
		SKUID:         "sku-1",
		SnapshotDate:  time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC),
		Price:         "1000",
		OriginalPrice: "1200",
		Currency:      "KRW",
	}); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	req := newTrackedItemRequest(t, http.MethodGet, "/v1/tracked-items/tracked-1/sku-price-trend?sku_id=sku-1&from=2026-04-01&to=2026-04-06", "tracked-1")

	status, body, err := controller.GetSKUPriceTrend(req)
	if err != nil {
		t.Fatalf("GetSKUPriceTrend() error = %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("GetSKUPriceTrend() status = %d, want 200", status)
	}

	resp, ok := body.(apiresponse.APIResponse[response.PriceTrendResponse])
	if !ok {
		t.Fatalf("GetSKUPriceTrend() body type = %T", body)
	}
	if resp.Data == nil {
		t.Fatal("GetSKUPriceTrend() returned nil data")
	}
	if len(resp.Data.Points) != 2 {
		t.Fatalf("GetSKUPriceTrend() points = %d, want 2", len(resp.Data.Points))
	}
	if resp.Data.Points[0].Date != "2026-04-04" {
		t.Fatalf("GetSKUPriceTrend() first date = %q, want 2026-04-04", resp.Data.Points[0].Date)
	}
	if resp.Data.Points[1].Date != "2026-04-05" {
		t.Fatalf("GetSKUPriceTrend() second date = %q, want 2026-04-05", resp.Data.Points[1].Date)
	}
	if resp.Data.Points[1].Currency != "KRW" {
		t.Fatalf("GetSKUPriceTrend() currency = %q, want KRW", resp.Data.Points[1].Currency)
	}
}

func TestGetDetail_UsesTodaySKUSnapshotPriceFirst(t *testing.T) {
	controller, _, snapshotRepo := newTestTrackedItemController(t, "KRW")

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	if err := snapshotRepo.Upsert(context.Background(), domainps.SKUPriceSnapshot{
		SKUID:         "sku-1",
		SnapshotDate:  yesterday,
		Price:         "15800",
		OriginalPrice: "16800",
		Currency:      "KRW",
	}); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if err := snapshotRepo.Upsert(context.Background(), domainps.SKUPriceSnapshot{
		SKUID:         "sku-1",
		SnapshotDate:  today,
		Price:         "15700",
		OriginalPrice: "16800",
		Currency:      "KRW",
	}); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	req := newTrackedItemRequest(t, http.MethodGet, "/v1/tracked-items/tracked-1", "tracked-1")

	status, body, err := controller.GetDetail(req)
	if err != nil {
		t.Fatalf("GetDetail() error = %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("GetDetail() status = %d, want 200", status)
	}

	resp, ok := body.(apiresponse.APIResponse[response.TrackedItemDetail])
	if !ok {
		t.Fatalf("GetDetail() body type = %T", body)
	}
	if resp.Data == nil {
		t.Fatal("GetDetail() returned nil data")
	}
	if resp.Data.CurrentPrice != "15700" {
		t.Fatalf("GetDetail() current_price = %q, want 15700", resp.Data.CurrentPrice)
	}
	if len(resp.Data.SKUs) == 0 || resp.Data.SKUs[0].CurrentPrice != "15700" {
		t.Fatalf("GetDetail() skus[0].current_price = %q, want 15700", firstSKUCurrentPrice(resp.Data.SKUs))
	}
}

func TestGetDetail_FallsBackToLatestSKUSnapshotPrice(t *testing.T) {
	controller, _, snapshotRepo := newTestTrackedItemController(t, "KRW")

	if err := snapshotRepo.Upsert(context.Background(), domainps.SKUPriceSnapshot{
		SKUID:         "sku-1",
		SnapshotDate:  time.Date(2026, 4, 4, 0, 0, 0, 0, time.UTC),
		Price:         "16000",
		OriginalPrice: "16800",
		Currency:      "KRW",
	}); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
	if err := snapshotRepo.Upsert(context.Background(), domainps.SKUPriceSnapshot{
		SKUID:         "sku-1",
		SnapshotDate:  time.Date(2026, 4, 9, 0, 0, 0, 0, time.UTC),
		Price:         "15650",
		OriginalPrice: "16800",
		Currency:      "KRW",
	}); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	req := newTrackedItemRequest(t, http.MethodGet, "/v1/tracked-items/tracked-1", "tracked-1")

	status, body, err := controller.GetDetail(req)
	if err != nil {
		t.Fatalf("GetDetail() error = %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("GetDetail() status = %d, want 200", status)
	}

	resp, ok := body.(apiresponse.APIResponse[response.TrackedItemDetail])
	if !ok {
		t.Fatalf("GetDetail() body type = %T", body)
	}
	if resp.Data == nil {
		t.Fatal("GetDetail() returned nil data")
	}
	if resp.Data.CurrentPrice != "15650" {
		t.Fatalf("GetDetail() current_price = %q, want 15650", resp.Data.CurrentPrice)
	}
	if len(resp.Data.SKUs) == 0 || resp.Data.SKUs[0].CurrentPrice != "15650" {
		t.Fatalf("GetDetail() skus[0].current_price = %q, want 15650", firstSKUCurrentPrice(resp.Data.SKUs))
	}
}

func firstSKUCurrentPrice(skus []response.ProductSKU) string {
	if len(skus) == 0 {
		return ""
	}
	return skus[0].CurrentPrice
}

func newTestTrackedItemController(t *testing.T, trackedItemCurrency string) (*Controller, *memoryskupricehistory.MemoryRepository, *memorypricesnapshot.SKUSnapshotMemoryRepository) {
	t.Helper()

	productRepo := memoryproduct.NewRepository()
	variantRepo := memoryproduct.NewVariantRepository()
	skuRepo := memoryproduct.NewSKURepository()
	trackedRepo := memorytrackeditem.NewRepository()
	historyRepo := memoryskupricehistory.NewRepository()
	snapshotRepo := memorypricesnapshot.NewSKUSnapshotRepository()
	productPriceHistoryRepo := newStubPriceHistoryRepository()

	productService := domainproduct.NewService(
		domainproduct.NewFinder(productRepo),
		domainproduct.NewWriter(productRepo),
		variantRepo,
		skuRepo,
		nil,
		nil,
		domainpricehistory.NewWriter(productPriceHistoryRepo),
		domainsph.NewWriter(historyRepo),
	)

	trackedService := domaintrackeditem.NewService(
		domaintrackeditem.NewFinder(trackedRepo),
		domaintrackeditem.NewWriter(trackedRepo),
		nil,
		nil,
		productService,
		nil,
	)

	if err := productRepo.Create(context.Background(), domainproduct.Product{
		ID:                "product-1",
		ExternalProductID: "external-1",
		OriginalURL:       "https://example.com/product/1",
		Title:             "Test Product",
	}); err != nil {
		t.Fatalf("productRepo.Create() error = %v", err)
	}
	if err := skuRepo.Create(context.Background(), domainproduct.SKU{
		ID:            "sku-1",
		ProductID:     "product-1",
		ExternalSKUID: "external-sku-1",
		Price:         "1000",
		OriginalPrice: "1200",
		Currency:      "KRW",
	}); err != nil {
		t.Fatalf("skuRepo.Create() error = %v", err)
	}
	if err := trackedRepo.Create(context.Background(), domaintrackeditem.TrackedItem{
		ID:        "tracked-1",
		UserID:    "",
		ProductID: "product-1",
		SKUID:     "sku-1",
		Currency:  trackedItemCurrency,
		CreatedAt: time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("trackedRepo.Create() error = %v", err)
	}

	controller := NewController(
		trackedService,
		domainsph.NewService(domainsph.NewFinder(historyRepo)),
		domainps.NewService(domainps.NewSKUSnapshotFinder(snapshotRepo)),
		nil,
	)
	return controller, historyRepo, snapshotRepo
}

func newTrackedItemRequest(t *testing.T, method string, target string, trackedItemID string) *http.Request {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), method, target, nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext() error = %v", err)
	}

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("trackedItemID", trackedItemID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeContext)
	return req.WithContext(ctx)
}

type stubPriceHistoryRepository struct{}

func newStubPriceHistoryRepository() stubPriceHistoryRepository {
	return stubPriceHistoryRepository{}
}

func (stubPriceHistoryRepository) Create(context.Context, domainpricehistory.PriceHistory) error {
	return nil
}

func (stubPriceHistoryRepository) ListByProductID(context.Context, string, string) ([]domainpricehistory.PriceHistory, error) {
	return nil, nil
}

type testAlertIDGenerator struct{}

func (testAlertIDGenerator) New() (string, error) {
	return "alert-1", nil
}

type testAlertClock struct{}

func (testAlertClock) Now() time.Time {
	return time.Unix(0, 0)
}
