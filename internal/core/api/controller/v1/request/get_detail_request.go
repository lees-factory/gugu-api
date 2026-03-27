package request

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type GetTrackedItemDetail struct {
	User          auth.RequestUser
	TrackedItemID string
}

func ParseGetTrackedItemDetail(r *http.Request) GetTrackedItemDetail {
	return GetTrackedItemDetail{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
	}
}
