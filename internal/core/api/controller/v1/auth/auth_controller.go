package auth

import (
	"net"
	stdhttp "net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
	supportauth "github.com/ljj/gugu-api/internal/support/auth"
)

type AuthController struct {
	authService *supportauth.Service
}

const defaultDeviceName = "unknown-device"

func NewAuthController(authService *supportauth.Service) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) RegisterRoutes(r chi.Router) {
	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/login/email", apiadvice.Wrap(c.LoginEmail))
		r.Post("/oauth/login", apiadvice.Wrap(c.LoginOAuth))
		r.Post("/refresh", apiadvice.Wrap(c.Refresh))
		r.Post("/logout", apiadvice.Wrap(c.Logout))
	})
}

func (c *AuthController) RegisterProtectedRoutes(r chi.Router) {
	r.Get("/v1/auth/sessions", apiadvice.Wrap(c.ListMySessions))
	r.Delete("/v1/auth/sessions/{sessionID}", apiadvice.Wrap(c.RevokeMySession))
}

func (c *AuthController) LoginEmail(r *stdhttp.Request) (int, any, error) {
	var req request.LoginEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	loginResult, err := c.authService.LoginEmail(r.Context(), supportauth.LoginEmailInput{
		Email:      req.Email,
		Password:   req.Password,
		UserAgent:  r.UserAgent(),
		ClientIP:   normalizedClientIP(r.RemoteAddr),
		DeviceName: normalizedDeviceName(r.Header.Get("X-Device-Name")),
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewLogin(*loginResult),
	), nil
}

func (c *AuthController) LoginOAuth(r *stdhttp.Request) (int, any, error) {
	var req request.LoginOAuth
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	loginResult, err := c.authService.LoginOAuth(r.Context(), supportauth.OAuthLoginInput{
		Provider:    supportauth.OAuthProvider(req.Provider),
		Subject:     req.Subject,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		UserAgent:   r.UserAgent(),
		ClientIP:    normalizedClientIP(r.RemoteAddr),
		DeviceName:  normalizedDeviceName(r.Header.Get("X-Device-Name")),
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewLogin(*loginResult),
	), nil
}

func (c *AuthController) Refresh(r *stdhttp.Request) (int, any, error) {
	var req request.RefreshToken
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	tokens, err := c.authService.RefreshTokens(r.Context(), supportauth.RefreshTokensInput{
		RefreshToken: req.RefreshToken,
		UserAgent:    r.UserAgent(),
		ClientIP:     normalizedClientIP(r.RemoteAddr),
		DeviceName:   normalizedDeviceName(req.DeviceName),
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewTokens(*tokens)), nil
}

func (c *AuthController) Logout(r *stdhttp.Request) (int, any, error) {
	var req request.Logout
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	if err := c.authService.Logout(r.Context(), supportauth.LogoutInput{RefreshToken: req.RefreshToken}); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}

func (c *AuthController) ListMySessions(r *stdhttp.Request) (int, any, error) {
	req := request.ParseListMyAuthSessions(r)

	sessions, err := c.authService.ListMyActiveSessions(r.Context(), req.User.ID)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(response.NewLoginSessions(sessions)), nil
}

func (c *AuthController) RevokeMySession(r *stdhttp.Request) (int, any, error) {
	req := request.ParseRevokeMyAuthSession(r)

	if err := c.authService.RevokeMySession(r.Context(), req.User.ID, req.SessionID); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}

func normalizedClientIP(remoteAddr string) string {
	trimmed := strings.TrimSpace(remoteAddr)
	if trimmed == "" {
		return ""
	}

	host, _, err := net.SplitHostPort(trimmed)
	if err == nil {
		if parsed := net.ParseIP(host); parsed != nil {
			return parsed.String()
		}
		return host
	}

	if parsed := net.ParseIP(trimmed); parsed != nil {
		return parsed.String()
	}

	return trimmed
}

func normalizedDeviceName(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return defaultDeviceName
	}
	return trimmed
}
