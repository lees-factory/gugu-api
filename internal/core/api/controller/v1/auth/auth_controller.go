package auth

import (
	stdhttp "net/http"

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

func (c *AuthController) LoginEmail(r *stdhttp.Request) (int, any, error) {
	var req request.LoginEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	loginResult, err := c.authService.LoginEmail(r.Context(), supportauth.LoginEmailInput{
		Email:      req.Email,
		Password:   req.Password,
		UserAgent:  r.UserAgent(),
		ClientIP:   r.RemoteAddr,
		DeviceName: r.Header.Get("X-Device-Name"),
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
		ClientIP:    r.RemoteAddr,
		DeviceName:  r.Header.Get("X-Device-Name"),
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
		ClientIP:     r.RemoteAddr,
		DeviceName:   req.DeviceName,
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
