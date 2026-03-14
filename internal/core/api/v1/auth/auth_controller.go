package auth

import (
	stdhttp "net/http"

	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	authrequest "github.com/ljj/gugu-api/internal/core/api/v1/auth/request"
	authresponse "github.com/ljj/gugu-api/internal/core/api/v1/auth/response"
	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

type Controller struct {
	authService *domainauth.Service
}

func NewController(authService *domainauth.Service) *Controller {
	return &Controller{authService: authService}
}

func (c *Controller) RegisterEmail(r *stdhttp.Request) (int, any, error) {
	var req authrequest.RegisterEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	result, err := c.authService.RegisterEmail(r.Context(), domainauth.RegisterEmailInput{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(
		authresponse.NewRegisterEmail(result.User, result.VerificationCode, result.VerificationDispatched),
	), nil
}

func (c *Controller) LoginEmail(r *stdhttp.Request) (int, any, error) {
	var req authrequest.LoginEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	loginResult, err := c.authService.LoginEmail(r.Context(), domainauth.LoginEmailInput{
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
		authresponse.NewLogin(*loginResult),
	), nil
}

func (c *Controller) VerifyEmail(r *stdhttp.Request) (int, any, error) {
	var req authrequest.VerifyEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	verifiedResult, err := c.authService.VerifyEmail(r.Context(), domainauth.VerifyEmailInput{Code: req.Code})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		authresponse.NewVerifyEmail(verifiedResult.User),
	), nil
}

func (c *Controller) LoginOAuth(r *stdhttp.Request) (int, any, error) {
	var req authrequest.LoginOAuth
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	loginResult, err := c.authService.LoginOAuth(r.Context(), domainauth.OAuthLoginInput{
		Provider:    domainauth.Provider(req.Provider),
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
		authresponse.NewLogin(*loginResult),
	), nil
}

func (c *Controller) Refresh(r *stdhttp.Request) (int, any, error) {
	var req authrequest.RefreshToken
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	tokens, err := c.authService.RefreshTokens(r.Context(), domainauth.RefreshTokensInput{
		RefreshToken: req.RefreshToken,
		UserAgent:    r.UserAgent(),
		ClientIP:     r.RemoteAddr,
		DeviceName:   req.DeviceName,
	})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(authresponse.NewTokens(*tokens)), nil
}

func (c *Controller) Logout(r *stdhttp.Request) (int, any, error) {
	var req authrequest.Logout
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	if err := c.authService.Logout(r.Context(), domainauth.LogoutInput{RefreshToken: req.RefreshToken}); err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.Success(), nil
}
