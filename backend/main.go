package main

import (
	"log"
	"urumi-backend/handlers"
	"urumi-backend/models"

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

	r := gin.Default()

	// CORS Setup (Simple for now)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Handlers
	storeHandler := handlers.NewStoreHandler(db)

	api := r.Group("/api")
	{
		api.GET("/stores", storeHandler.ListStores)
		api.POST("/stores", storeHandler.CreateStore)
		api.DELETE("/stores/:id", storeHandler.DeleteStore)
	}

	r.Run(":8080")
}
