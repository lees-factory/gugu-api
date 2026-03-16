package product

import (
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	productresponse "github.com/ljj/gugu-api/internal/core/api/v1/product/response"
	productdetail "github.com/ljj/gugu-api/internal/core/domain/productdetail"
)

type Controller struct {
	service *productdetail.Service
}

func NewController(service *productdetail.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetDetail(r *stdhttp.Request) (int, any, error) {
	productID := strings.TrimSpace(chi.URLParam(r, "productID"))
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))

	result, err := c.service.Get(r.Context(), productID, userID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(productresponse.NewProductDetail(*result)), nil
}
