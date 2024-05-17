package router

import (
	"github.com/gin-contrib/cors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	spec "github.com/khedhrije/podcaster-search-api/deployments/swagger"
	"github.com/khedhrije/podcaster-search-api/internal/configuration"
	"github.com/khedhrije/podcaster-search-api/internal/ui/gin/handlers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// CreateRouter sets up and returns a new Gin router with the defined routes.
func CreateRouter(handler handlers.Search) *gin.Engine {
	// Initialize a new Gin router without any middleware by default.
	r := gin.New()

	// Customize CORS configuration if needed
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"}, // Change this to specific domains if needed
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}

	// Apply the CORS middleware with the custom configuration
	r.Use(cors.New(corsConfig))

	// Health check route.
	r.GET("/health", health())

	// Configure Swagger documentation URL based on the environment.
	configureSwagger()

	// Set up the route for Swagger documentation.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Define private routes that require authentication.
	private := r.Group("/private")
	private.Use(TokenValidatorMiddleware())
	{
		// Routes for managing search operations.
		search := private.Group("/search")
		{
			search.GET("/programs/", handler.Programs())
			search.GET("/programs/:id", handler.ProgramByID())
		}
	}

	// Return the configured router.
	return r
}

// configureSwagger configures the Swagger documentation URL based on the environment.
func configureSwagger() {
	if configuration.Config.Env == "dev" {
		spec.SwaggerInfo.Host = configuration.Config.DocsAddress
	} else {
		spec.SwaggerInfo.Host = removeHTTPS(configuration.Config.DocsAddress)
	}
}

// removeHTTPS removes the "https://" prefix from a URL string.
func removeHTTPS(url string) string {
	if strings.HasPrefix(url, "https://") {
		return strings.TrimPrefix(url, "https://")
	}
	return url
}
