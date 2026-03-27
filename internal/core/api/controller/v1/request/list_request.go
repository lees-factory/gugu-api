package request

import (
	"net/http"

	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type ListTrackedItems struct {
	User auth.RequestUser
}

func ParseListTrackedItems(r *http.Request) ListTrackedItems {
	return ListTrackedItems{
		User: auth.RequestUserFrom(r.Context()),
	}
}
