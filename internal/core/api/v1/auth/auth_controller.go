package auth

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"

	authrequest "github.com/ljj/gugu-api/internal/core/api/v1/auth/request"
	authresponse "github.com/ljj/gugu-api/internal/core/api/v1/auth/response"
	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
)

type Controller struct {
	service *domainauth.Service
}

func NewController(service *domainauth.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) RegisterEmail(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req authrequest.RegisterEmail
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid request body")
		return
	}

	result, err := c.service.RegisterEmail(r.Context(), domainauth.RegisterEmailInput{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		c.writeAuthError(w, err)
		return
	}

	writeJSON(w, stdhttp.StatusCreated, authresponse.NewRegisterEmail(*result))
}

func (c *Controller) LoginEmail(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req authrequest.LoginEmail
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid request body")
		return
	}

	result, err := c.service.LoginEmail(r.Context(), domainauth.LoginEmailInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.writeAuthError(w, err)
		return
	}

	writeJSON(w, stdhttp.StatusOK, authresponse.NewLogin(*result))
}

func (c *Controller) VerifyEmail(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req authrequest.VerifyEmail
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid request body")
		return
	}

	result, err := c.service.VerifyEmail(r.Context(), domainauth.VerifyEmailInput{Token: req.Token})
	if err != nil {
		c.writeAuthError(w, err)
		return
	}

	writeJSON(w, stdhttp.StatusOK, authresponse.NewVerifyEmail(*result))
}

func (c *Controller) LoginOAuth(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req authrequest.LoginOAuth
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, "invalid request body")
		return
	}

	result, err := c.service.LoginOAuth(r.Context(), domainauth.OAuthLoginInput{
		Provider:    domainauth.Provider(req.Provider),
		Subject:     req.Subject,
		Email:       req.Email,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		c.writeAuthError(w, err)
		return
	}

	writeJSON(w, stdhttp.StatusOK, authresponse.NewLogin(*result))
}

func (c *Controller) writeAuthError(w stdhttp.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domainauth.ErrEmailAlreadyExists):
		writeError(w, stdhttp.StatusConflict, err.Error())
	case errors.Is(err, domainauth.ErrInvalidCredentials):
		writeError(w, stdhttp.StatusUnauthorized, err.Error())
	case errors.Is(err, domainauth.ErrEmailNotVerified):
		writeError(w, stdhttp.StatusForbidden, err.Error())
	case errors.Is(err, domainauth.ErrVerificationNotFound), errors.Is(err, domainauth.ErrOAuthProviderInvalid):
		writeError(w, stdhttp.StatusBadRequest, err.Error())
	default:
		writeError(w, stdhttp.StatusInternalServerError, "internal server error")
	}
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w stdhttp.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
