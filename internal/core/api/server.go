package api

import (
	"database/sql"
	"fmt"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	apiauth "github.com/ljj/gugu-api/internal/core/api/controller/v1/auth"
	apiintegration "github.com/ljj/gugu-api/internal/core/api/controller/v1/integration"
	apiproduct "github.com/ljj/gugu-api/internal/core/api/controller/v1/product"
	apitrackeditem "github.com/ljj/gugu-api/internal/core/api/controller/v1/trackeditem"
	"github.com/ljj/gugu-api/internal/core/api/admin"
	"github.com/ljj/gugu-api/internal/core/support/auth"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/security"
)

type Server struct {
	router *chi.Mux
}

func NewServer(cfg config.Config, db *sql.DB) (*Server, error) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(CORSMiddleware(cfg.CORSAllowedOrigins))

	registerHealthRoute(router)

	authControllers, err := apiauth.Wire(cfg, db)
	if err != nil {
		return nil, fmt.Errorf("wire auth: %w", err)
	}
	authControllers.Auth.RegisterRoutes(router)
	authControllers.User.RegisterRoutes(router)

	aliExpressController, aliExpressTokenStore, err := apiintegration.Wire(cfg, db)
	if err != nil {
		return nil, fmt.Errorf("wire aliexpress integration: %w", err)
	}
	aliExpressController.RegisterRoutes(router)

	trackedItemController, trackedItemService, productService, err := apitrackeditem.Wire(cfg, db, aliExpressTokenStore)
	if err != nil {
		return nil, fmt.Errorf("wire tracked item: %w", err)
	}

	jwtIssuer := security.NewJWTTokenIssuer(cfg.JWTSecret, cfg.JWTIssuer)
	router.Group(func(r chi.Router) {
		r.Use(auth.UserArgumentResolver(jwtIssuer))
		trackedItemController.RegisterRoutes(r)
	})

	productController := apiproduct.Wire(db, productService, trackedItemService)
	productController.RegisterRoutes(router)

	admin.RegisterBatchRoutes(router, cfg, db, aliExpressTokenStore, productService)

	return &Server{router: router}, nil
}

func (s *Server) Handler() stdhttp.Handler {
	return s.router
}

func registerHealthRoute(router chi.Router) {
	router.Get("/health", apiadvice.Wrap(func(_ *stdhttp.Request) (int, any, error) {
		return stdhttp.StatusOK, apiresponse.SuccessWithData(map[string]string{"status": "ok"}), nil
	}))
}
