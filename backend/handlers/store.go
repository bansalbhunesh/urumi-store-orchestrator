package handlers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"urumi-backend/models"
	"urumi-backend/orchestrator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreHandler struct {
	DB *gorm.DB
}

func NewStoreHandler(db *gorm.DB) *StoreHandler {
	return &StoreHandler{DB: db}
}

func (h *StoreHandler) ListStores(c *gin.Context) {
	var stores []models.Store
	h.DB.Find(&stores)
	c.JSON(http.StatusOK, stores)
}

func (h *StoreHandler) CreateStore(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Validate store name
	if len(strings.TrimSpace(input.Name)) < 2 || len(input.Name) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Store name must be between 2 and 50 characters"})
		return
	}

	// Sanitize store name
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !nameRegex.MatchString(input.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Store name can only contain letters, numbers, spaces, hyphens, and underscores"})
		return
	}

	// Validate store type
	if input.Type != "woocommerce" && input.Type != "medusa" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Store type must be either 'woocommerce' or 'medusa'"})
		return
	}

	storeID := uuid.New().String()
	namespace := "store-" + storeID[:8]

	domainSuffix := os.Getenv("DOMAIN_SUFFIX")
	if domainSuffix == "" {
		domainSuffix = "localhost"
	}
	store := models.Store{
		ID:        storeID,
		Name:      strings.TrimSpace(input.Name),
		Type:      input.Type,
		Status:    "Provisioning",
		Namespace: namespace,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		URL:       "http://" + namespace + "." + domainSuffix,
	}

	if err := h.DB.Create(&store).Error; err != nil {
		log.Printf("Failed to create store record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store record"})
		return
	}

	// Trigger async provisioning
	go func(s models.Store) {
		log.Printf("Starting provisioning for store %s (%s)", s.ID, s.Name)
		err := orchestrator.ProvisionStore(s)
		status := "Ready"
		errorMessage := (*string)(nil)
		if err != nil {
			status = "Failed"
			errStr := err.Error()
			errorMessage = &errStr
			log.Printf("Failed to provision store %s: %v", s.ID, err)
		} else {
			log.Printf("Successfully provisioned store %s", s.ID)
		}
		
		updateErr := h.DB.Model(&s).Updates(map[string]interface{}{
			"status":         status,
			"error_message":  errorMessage,
			"updated_at":     time.Now(),
		}).Error
		
		if updateErr != nil {
			log.Printf("Failed to update store status for %s: %v", s.ID, updateErr)
		}
	}(store)

	c.JSON(http.StatusAccepted, store)
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id := c.Param("id")
	var store models.Store
	if result := h.DB.First(&store, "id = ?", id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		} else {
			log.Printf("Database error when fetching store %s: %v", id, result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check if store is already being deleted
	if store.Status == "Deleting" {
		c.JSON(http.StatusConflict, gin.H{"error": "Store is already being deleted"})
		return
	}

	// Mark as deleting immediately for UI feedback
	if err := h.DB.Model(&store).Updates(map[string]interface{}{
		"status":     "Deleting",
		"updated_at": time.Now(),
	}).Error; err != nil {
		log.Printf("Failed to mark store %s as deleting: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store status"})
		return
	}

	// Trigger async deletion
	go func(s models.Store) {
		log.Printf("Starting deletion for store %s (%s)", s.ID, s.Name)
		err := orchestrator.DeleteStore(s)
		if err != nil {
			log.Printf("Failed to delete store %s: %v", s.ID, err)
			// Mark as failed deletion
			h.DB.Model(&s).Updates(map[string]interface{}{
				"status":        "DeletionFailed",
				"error_message": &[]string{err.Error()}[0],
				"updated_at":    time.Now(),
			})
		} else {
			log.Printf("Successfully deleted store %s", s.ID)
			// Remove from database only after successful deletion
			h.DB.Delete(&s)
		}
	}(store)

	c.JSON(http.StatusOK, gin.H{"message": "Store deletion started"})
}

func (h *StoreHandler) CheckStoreHealth(c *gin.Context) {
	id := c.Param("id")
	var store models.Store
	if result := h.DB.First(&store, "id = ?", id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		} else {
			log.Printf("Database error when fetching store %s: %v", id, result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Perform health check
	healthy, err := orchestrator.CheckStoreHealth(store)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"healthy": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"healthy": healthy,
		"status":  store.Status,
	})
}
