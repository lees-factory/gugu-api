package discover

import (
	stdhttp "net/http"
	"strconv"

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

	products, err := c.productService.ListByCollectionSource(r.Context(), "HOT_PRODUCT_QUERY", size, offset)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewHotProductList(products)), nil
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
