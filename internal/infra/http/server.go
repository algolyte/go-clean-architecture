package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zahidhasanpapon/go-clean-architecture/config"
	"github.com/zahidhasanpapon/go-clean-architecture/internal/infra/http/middleware"
	"github.com/zahidhasanpapon/go-clean-architecture/pkg/logger"
	"net/http"
)

type Server struct {
	router     *gin.Engine
	config     *config.Config
	logger     logger.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.Config, l logger.Logger) (*Server, error) {
	// Set Gin mode based on environment
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	s := &Server{
		router: router,
		config: cfg,
		logger: l,
	}

	// Setup middleware and routes
	s.setupMiddleware()
	s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	return s, nil
}

func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Custom middleware
	s.router.Use(middleware.RequestID())
	s.router.Use(middleware.Logger(s.logger))
	s.router.Use(middleware.CORS(s.config))
	s.router.Use(middleware.ErrorHandler())
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}
}

func (s *Server) Start() error {
	s.logger.Info(context.Background(), "Starting server on port %s")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info(ctx, "Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
