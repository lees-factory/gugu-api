package discover

import (
	stdhttp "net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

type Controller struct {
	productService *domainproduct.Service
}

func NewController(productService *domainproduct.Service) *Controller {
	return &Controller{productService: productService}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Route("/v1/discover", func(r chi.Router) {
		r.Get("/hot-products", apiadvice.Wrap(c.ListHotProducts))
	})
}

func (c *Controller) ListHotProducts(r *stdhttp.Request) (int, any, error) {
	page, size := parsePagination(r)
	offset := (page - 1) * size
	language := parseHotProductLanguage(r)
	collectionSource := hotProductCollectionSource(language)

	products, err := c.productService.ListByCollectionSource(r.Context(), collectionSource, size, offset)
	if err != nil {
		return 0, nil, err
	}
	if len(products) == 0 && language == "KO" {
		products, err = c.productService.ListByCollectionSource(r.Context(), "HOT_PRODUCT_QUERY", size, offset)
		if err != nil {
			return 0, nil, err
		}
	}

	keys := make([]domainproduct.VariantLookupKey, 0, len(products))
	for _, product := range products {
		keys = append(keys, domainproduct.VariantLookupKey{
			ProductID: product.ID,
			Language:  language,
			Currency:  hotProductCurrency(language),
		})
	}
	foundVariants, err := c.productService.FindVariants(r.Context(), keys)
	if err != nil {
		return 0, nil, err
	}
	variants := make(map[string]*domainproduct.Variant, len(foundVariants))
	for i := range foundVariants {
		variant := foundVariants[i]
		variants[variant.ProductID] = &variant
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewHotProductList(products, variants, hotProductCurrency(language))), nil
}

func parsePagination(r *stdhttp.Request) (page int, size int) {
	page = 1
	size = defaultPageSize

	if v, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(r.URL.Query().Get("size")); err == nil && v > 0 {
		size = v
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return page, size
}

func hotProductCollectionSource(targetLanguage string) string {
	return "HOT_PRODUCT_QUERY:" + strings.ToUpper(strings.TrimSpace(targetLanguage))
}

func parseHotProductLanguage(r *stdhttp.Request) string {
	language := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("language")))
	if language != "" {
		return language
	}

	language = strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("target_language")))
	if language != "" {
		return language
	}

	return "KO"
}

func hotProductCurrency(language string) string {
	switch strings.ToUpper(strings.TrimSpace(language)) {
	case "KO":
		return "KRW"
	default:
		return "USD"
	}
}
