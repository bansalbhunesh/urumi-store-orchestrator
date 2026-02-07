package models

import (
	"time"
)

type Store struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "woocommerce" or "medusa"
	Status    string    `json:"status"` // Provisioning, Ready, Failed
	URL       string    `json:"url"`
	Namespace string    `json:"namespace"`
	CreatedAt time.Time `json:"created_at"`
}
