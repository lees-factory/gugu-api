package auth

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	authrequest "github.com/ljj/gugu-api/internal/core/api/v1/auth/request"
	authresponse "github.com/ljj/gugu-api/internal/core/api/v1/auth/response"
	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
)

type UserController struct {
	userService *domainuser.Service
}

func NewUserController(userService *domainuser.Service) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) RegisterRoutes(r chi.Router) {
	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register/email", apiadvice.Wrap(c.RegisterEmail))
		r.Post("/verify-email", apiadvice.Wrap(c.VerifyEmail))
	})
}

func (c *UserController) RegisterEmail(r *stdhttp.Request) (int, any, error) {
	var req authrequest.RegisterEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	result, err := c.userService.RegisterEmail(r.Context(), domainuser.RegisterEmailInput{
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

func (c *UserController) VerifyEmail(r *stdhttp.Request) (int, any, error) {
	var req authrequest.VerifyEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	verifiedResult, err := c.userService.VerifyEmail(r.Context(), domainuser.VerifyEmailInput{Code: req.Code})
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		authresponse.NewVerifyEmail(verifiedResult.User),
	), nil
}
