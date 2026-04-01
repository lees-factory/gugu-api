package api

import (
	"database/sql"
	"fmt"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	apiauth "github.com/ljj/gugu-api/internal/core/api/controller/v1/auth"
	apidiscover "github.com/ljj/gugu-api/internal/core/api/controller/v1/discover"
	apiintegration "github.com/ljj/gugu-api/internal/core/api/controller/v1/integration"
	apipricealert "github.com/ljj/gugu-api/internal/core/api/controller/v1/pricealert"
	apiproduct "github.com/ljj/gugu-api/internal/core/api/controller/v1/product"
	apitrackeditem "github.com/ljj/gugu-api/internal/core/api/controller/v1/trackeditem"
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
	registerOpenAPIRoute(router)

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

	priceAlertController := apipricealert.Wire(db)

	jwtIssuer := security.NewJWTTokenIssuer(cfg.JWTSecret, cfg.JWTIssuer)
	router.Group(func(r chi.Router) {
		r.Use(auth.UserArgumentResolver(jwtIssuer))
		trackedItemController.RegisterRoutes(r)
		priceAlertController.RegisterRoutes(r)
	})

	productController := apiproduct.Wire(db, productService, trackedItemService)
	productController.RegisterRoutes(router)

	discoverController := apidiscover.NewController(productService)
	discoverController.RegisterRoutes(router)

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

func registerOpenAPIRoute(router chi.Router) {
	router.Get("/openapi.yml", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		stdhttp.ServeFile(w, r, "openapi.yml")
	})
}
