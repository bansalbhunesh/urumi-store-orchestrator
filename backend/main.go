package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"urumi-backend/handlers"
	"urumi-backend/middleware"
	"urumi-backend/models"
	"urumi-backend/orchestrator"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize Database
	db, err := gorm.Open(sqlite.Open("stores.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Store{})

	// Start background reconciliation
	go startReconciliationService(db)

	// Initialize rate limiter (20 requests per minute, burst of 40) - increased for demo
	rateLimiter := middleware.NewRateLimiter(20, 40)
	go rateLimiter.CleanupExpiredClients()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Add security middlewares
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.TimeoutMiddleware(30 * time.Second))
	r.Use(middleware.ValidateContentType())
	r.Use(gin.Recovery())

	// Add CORS middleware with secure configuration
	corsConfig := middleware.DefaultCORSConfig()
	r.Use(middleware.CORSMiddleware(corsConfig))

	// Add rate limiting
	r.Use(rateLimiter.Middleware())

	// Request logging
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
	}))

	// Handlers
	storeHandler := handlers.NewStoreHandler(db)

	api := r.Group("/api")
	{
		api.GET("/stores", storeHandler.ListStores)
		api.POST("/stores", storeHandler.CreateStore)
		api.DELETE("/stores/:id", storeHandler.DeleteStore)
		api.GET("/stores/:id/health", storeHandler.CheckStoreHealth)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	log.Println("Starting Urumi Backend Server on :8080")
	log.Println("Security features enabled: CORS, Rate Limiting, Security Headers")
	r.Run(":8080")
}

// startReconciliationService runs background checks to ensure store status is accurate
func startReconciliationService(db *gorm.DB) {
	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	log.Println("Starting store reconciliation service")

	for range ticker.C {
		var stores []models.Store
		if err := db.Find(&stores).Error; err != nil {
			log.Printf("Failed to fetch stores for reconciliation: %v", err)
			continue
		}

		for _, store := range stores {
			// Skip stores that are being deleted
			if store.Status == "Deleting" || store.Status == "DeletionFailed" {
				continue
			}

			// Reconcile store status
			if err := orchestrator.ReconcileStoreStatus(store, db); err != nil {
				log.Printf("Failed to reconcile store %s: %v", store.ID, err)
			}
		}
	}
}
