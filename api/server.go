package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pranaykumar2/steg-go/api/docs"
	"github.com/pranaykumar2/steg-go/api/handlers"
	"github.com/pranaykumar2/steg-go/api/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	port   string
}

// NewServer creates a new API server
func NewServer() *Server {
	// Set port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set default host for Swagger docs
	host := os.Getenv("API_HOST")
	if host == "" {
		host = "localhost:" + port
	}
	docs.SwaggerInfo.Host = host

	// Create server instance
	server := &Server{
		router: gin.Default(),
		port:   port,
	}

	// Initialize routes and middleware
	server.setupMiddleware()
	server.setupRoutes()

	return server
}

// setupMiddleware configures all middleware for the server
func (s *Server) setupMiddleware() {
	// Add CORS middleware
	s.router.Use(middleware.CORS())

	// Add security headers
	s.router.Use(middleware.SecurityHeaders())

	// Add rate limiting - 100 requests per minute
	s.router.Use(middleware.RateLimit(100, 1*time.Minute))

	// File size limits - 10MB max upload
	s.router.MaxMultipartMemory = 50 << 20
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// API version group
	v1 := s.router.Group("/api")
	{
		// Health check endpoint
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"version": "1.0.0",
				"time":    time.Now().UTC().Format(time.RFC3339),
			})
		})

		// Steganography endpoints
		v1.POST("/hide", handlers.HideText)
		v1.POST("/hideFile", handlers.HideFile)
		v1.POST("/extract", handlers.Extract)
		v1.POST("/metadata", handlers.AnalyzeMetadata)

		// File serving endpoint
		v1.GET("/files/:filename", handlers.ServeFile)
	}

	// Swagger documentation endpoint
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static files for the Web UI
	s.router.Static("/static", "./web/static")

	// Handle Web UI routes
	s.setupWebUIRoutes()
}

// setupWebUIRoutes configures routes for the web interface
func (s *Server) setupWebUIRoutes() {
	// Serve index.html at root
	s.router.GET("/", func(c *gin.Context) {
		c.File("./web/templates/index.html")
	})

	// For any route not matched by API or static files,
	// serve index.html to support SPA routing
	s.router.NoRoute(func(c *gin.Context) {
		// If it's an API request that wasn't matched, return 404 JSON
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[0:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "API endpoint not found",
			})
			return
		}

		// If not an API request, serve the SPA
		c.File("./web/templates/index.html")
	})
}

// Start begins listening for requests
func (s *Server) Start() error {
	log.Printf("Starting StegGo server on port %s", s.port)
	log.Printf("Current time: %s", time.Now().UTC().Format("2006-01-02 15:04:05"))
	log.Printf("Web UI available at http://localhost:%s", s.port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", s.port)
	return s.router.Run(fmt.Sprintf(":%s", s.port))
}
