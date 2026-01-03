package server

import (
	"github.com/programmerjide/ecommerce/internal/handler"
	"github.com/programmerjide/ecommerce/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/programmerjide/ecommerce/internal/config"
	"github.com/programmerjide/ecommerce/internal/middleware"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
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

	authService := service.NewAuthService(s.db, s.config)
	userService := service.NewUserService(s.db)
	productService := service.NewProductService(s.db)

	authHandler := handler.NewAuthHandler(authService, *s.logger)
	userHandler := handler.NewUserHandler(userService, *s.logger)
	productHandler := handler.NewProductHandler(productService, *s.logger)

	api := router.Group("/api/v1") // API v1 routes
	{
		// Public routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(&s.config.JWT))
		{
			users := protected.Group("/users")
			{
				// User profile routes
				userRoutes := users
				userRoutes.GET("/profile", userHandler.GetProfile)
				userRoutes.PUT("/profile", userHandler.UpdateProfile)
			}

			categories := protected.Group("/categories")
			{
				categoryRoutes := categories
				categoryRoutes.POST("/", middleware.AdminMiddleware(), productHandler.CreateCategory)
				categoryRoutes.GET("/", productHandler.GetCategories)
				categoryRoutes.PUT("/:id", middleware.AdminMiddleware(), productHandler.UpdateCategory)
				categoryRoutes.DELETE("/:id", middleware.AdminMiddleware(), productHandler.DeleteCategory)
			}

			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", middleware.AdminMiddleware(), productHandler.CreateProduct)      // No ()
				productRoutes.GET("/", productHandler.GetProducts)                                       // No ()
				productRoutes.GET("/:id", productHandler.GetProduct)                                     // No () - FIXED
				productRoutes.PUT("/:id", middleware.AdminMiddleware(), productHandler.UpdateProduct)    // No ()
				productRoutes.DELETE("/:id", middleware.AdminMiddleware(), productHandler.DeleteProduct) // No ()
				productRoutes.GET("/search", productHandler.SearchProducts)                              // Changed from POST, moved before /:id
			}
		}
	}

	return router
}

func (s *Server) healthCheckHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
