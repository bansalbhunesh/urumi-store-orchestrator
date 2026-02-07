package handlers

import (
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storeID := uuid.New().String()
	namespace := "store-" + storeID[:8]
	
	store := models.Store{
		ID:        storeID,
		Name:      input.Name,
		Type:      input.Type,
		Status:    "Provisioning",
		Namespace: namespace,
		CreatedAt: time.Now(),
		URL:       "http://" + namespace + ".localhost", // Simplistic for now
	}

	h.DB.Create(&store)

	// Trigger async provisioning
	go func(s models.Store) {
		err := orchestrator.ProvisionStore(s)
		status := "Ready"
		if err != nil {
			status = "Failed"
		}
		h.DB.Model(&s).Update("status", status)
	}(store)

	c.JSON(http.StatusAccepted, store)
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id := c.Param("id")
	var store models.Store
	if result := h.DB.First(&store, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	// Trigger async deletion
	go func(s models.Store) {
		orchestrator.DeleteStore(s)
		h.DB.Delete(&s) // Ideally delete after confirmation, but for immediate UI feedback:
	}(store)
	
	// We delete from DB immediately for UI responsiveness in this simple demo, 
	// or we could mark as "Deleting"
	h.DB.Delete(&store)

	c.JSON(http.StatusOK, gin.H{"message": "Store deletion started"})
}
