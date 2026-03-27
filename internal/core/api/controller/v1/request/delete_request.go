package request

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type DeleteTrackedItem struct {
	User          auth.RequestUser
	TrackedItemID string
}

func ParseDeleteTrackedItem(r *http.Request) DeleteTrackedItem {
	return DeleteTrackedItem{
		User:          auth.RequestUserFrom(r.Context()),
		TrackedItemID: strings.TrimSpace(chi.URLParam(r, "trackedItemID")),
	}
}
