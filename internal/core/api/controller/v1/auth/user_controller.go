package auth

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	"github.com/ljj/gugu-api/internal/core/api/controller/v1/response"
	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserController struct {
	userService    *domainuser.Service
	passwordHasher PasswordHasher
}

func NewUserController(userService *domainuser.Service, passwordHasher PasswordHasher) *UserController {
	return &UserController{userService: userService, passwordHasher: passwordHasher}
}

func (c *UserController) RegisterRoutes(r chi.Router) {
	r.Post("/v1/auth/register/email", apiadvice.Wrap(c.RegisterEmail))
	r.Post("/v1/auth/verify-email", apiadvice.Wrap(c.VerifyEmail))
}

func (c *UserController) RegisterEmail(r *stdhttp.Request) (int, any, error) {
	var req request.RegisterEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	hash, err := c.passwordHasher.Hash(req.Password)
	if err != nil {
		return 0, nil, err
	}

	user, err := c.userService.Create(r.Context(), req.ToNewUser(hash))
	if err != nil {
		return 0, nil, err
	}

	code, err := c.userService.SendVerification(r.Context(), user.ID, user.Email)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusCreated, apiresponse.SuccessWithData(
		response.NewRegisterEmail(*user, code, true),
	), nil
}

func (c *UserController) VerifyEmail(r *stdhttp.Request) (int, any, error) {
	var req request.VerifyEmail
	if err := apiadvice.DecodeJSON(r, &req); err != nil {
		return 0, nil, err
	}

	user, err := c.userService.VerifyEmail(r.Context(), req.Code)
	if err != nil {
		return 0, nil, err
	}

	return stdhttp.StatusOK, apiresponse.SuccessWithData(
		response.NewVerifyEmail(*user),
	), nil
}
