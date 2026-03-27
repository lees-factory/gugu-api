package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ljj/gugu-api/internal/support/security"
)

type contextKey string

const requestUserContextKey contextKey = "requestUser"

type RequestUser struct {
	ID string
}

func UserArgumentResolver(jwtIssuer security.JWTTokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				writeAuthError(w, http.StatusUnauthorized, "E401", "Authorization header is required")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := jwtIssuer.VerifyAccessToken(token)
			if err != nil {
				writeAuthError(w, http.StatusUnauthorized, "E401", "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), requestUserContextKey, RequestUser{ID: userID})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequestUserFrom(ctx context.Context) RequestUser {
	if v, ok := ctx.Value(requestUserContextKey).(RequestUser); ok {
		return v
	}
	return RequestUser{}
}

func writeAuthError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"result": "ERROR",
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
