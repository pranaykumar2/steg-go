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

type Server struct {
	router *gin.Engine
	port   string
}

func NewServer() *Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("API_HOST")
	if host == "" {
		host = "localhost:" + port
	}
	docs.SwaggerInfo.Host = host

	server := &Server{
		router: gin.Default(),
		port:   port,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.CORS())

	s.router.Use(middleware.SecurityHeaders())

	s.router.Use(middleware.RateLimit(100, 1*time.Minute))

	s.router.MaxMultipartMemory = 50 << 20
}

func (s *Server) setupRoutes() {
	v1 := s.router.Group("/api")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"version": "1.0.0",
				"time":    time.Now().UTC().Format(time.RFC3339),
			})
		})

		v1.POST("/hide", handlers.HideText)
		v1.POST("/hideFile", handlers.HideFile)
		v1.POST("/extract", handlers.Extract)
		v1.POST("/metadata", handlers.AnalyzeMetadata)

		// File serving endpoint
		v1.GET("/files/:filename", handlers.ServeFile)
	}
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.router.Static("/static", "./web/static")

	s.setupWebUIRoutes()
}

func (s *Server) setupWebUIRoutes() {
	// Serve index.html at root
	s.router.GET("/", func(c *gin.Context) {
		c.File("./web/templates/index.html")
	})

	s.router.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[0:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "API endpoint not found",
			})
			return
		}

		c.File("./web/templates/index.html")
	})
}

func (s *Server) Start() error {
	log.Printf("Starting StegGo server on port %s", s.port)
	log.Printf("Current time: %s", time.Now().UTC().Format("2006-01-02 15:04:05"))
	log.Printf("Web UI available at http://localhost:%s", s.port)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", s.port)
	return s.router.Run(fmt.Sprintf(":%s", s.port))
}
