package request

import (
	"net/http"
	"strconv"

	"github.com/ljj/gugu-api/internal/core/support/auth"
	"github.com/ljj/gugu-api/internal/core/support/page"
)

type ListTrackedItems struct {
	User   auth.RequestUser
	Cursor page.CursorRequest
}

func ParseListTrackedItems(r *http.Request) ListTrackedItems {
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	return ListTrackedItems{
		User: auth.RequestUserFrom(r.Context()),
		Cursor: page.CursorRequest{
			Cursor: r.URL.Query().Get("cursor"),
			Size:   size,
		},
	}
}
