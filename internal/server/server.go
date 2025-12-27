package server

import (
	"github.com/gin-gonic/gin"
	"github.com/programmerjide/ecommerce/internal/config"
	"github.com/programmerjide/ecommerce/internal/middleware"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewServer(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	return &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	// Placeholder for setting up routes
	router := gin.New()

	// Add middleware, routes, etc.
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CORS(s.config))

	// Add routes here
	router.GET("/health", s.healthCheckHandler)

	return router
}

func (s *Server) healthCheckHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
