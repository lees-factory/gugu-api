package request

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ljj/gugu-api/internal/core/support/auth"
)

type ListMyAuthSessions struct {
	User auth.RequestUser
}

func ParseListMyAuthSessions(r *http.Request) ListMyAuthSessions {
	return ListMyAuthSessions{User: auth.RequestUserFrom(r.Context())}
}

type RevokeMyAuthSession struct {
	User      auth.RequestUser
	SessionID string
}

func ParseRevokeMyAuthSession(r *http.Request) RevokeMyAuthSession {
	return RevokeMyAuthSession{
		User:      auth.RequestUserFrom(r.Context()),
		SessionID: strings.TrimSpace(chi.URLParam(r, "sessionID")),
	}
}
