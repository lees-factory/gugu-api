package request

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type GetProductDetail struct {
	User      auth.RequestUser
	ProductID string
}

func ParseGetProductDetail(r *http.Request) GetProductDetail {
	return GetProductDetail{
		User:      auth.RequestUserFrom(r.Context()),
		ProductID: strings.TrimSpace(chi.URLParam(r, "productID")),
	}
}

type ListProductSKUs struct {
	ProductID string
}

func ParseListProductSKUs(r *http.Request) ListProductSKUs {
	return ListProductSKUs{
		ProductID: strings.TrimSpace(chi.URLParam(r, "productID")),
	}
}
