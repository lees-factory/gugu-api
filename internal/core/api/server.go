package api

import (
	"database/sql"
	"fmt"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	"github.com/ljj/gugu-api/internal/support/config"
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

	if err := registerAuthModule(router, cfg, db); err != nil {
		return nil, fmt.Errorf("register auth module: %w", err)
	}
	if err := registerAliExpressIntegrationModule(router, cfg, db); err != nil {
		return nil, fmt.Errorf("register aliexpress integration module: %w", err)
	}
	if err := registerProductModule(router, cfg, db); err != nil {
		return nil, fmt.Errorf("register product module: %w", err)
	}

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
