package orchestrator

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"urumi-backend/models"
)

// CheckStoreHealth performs a health check on a provisioned store
func CheckStoreHealth(store models.Store) (bool, error) {
	// For WooCommerce stores, check if the WordPress site is responding
	if store.Type == "woocommerce" {
		return checkWooCommerceHealth(store)
	}
	
	// For Medusa stores, check if the API is responding
	if store.Type == "medusa" {
		return checkMedusaHealth(store)
	}
	
	return false, fmt.Errorf("unknown store type: %s", store.Type)
}

func checkWooCommerceHealth(store models.Store) (bool, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Try to access the WordPress site
	resp, err := client.Get(store.URL)
	if err != nil {
		log.Printf("Health check failed for %s: %v", store.URL, err)
		return false, err
	}
	defer resp.Body.Close()
	
	// Check if we get a successful response
	if resp.StatusCode == http.StatusOK {
		log.Printf("Health check passed for %s", store.URL)
		return true, nil
	}
	
	log.Printf("Health check failed for %s - status: %d", store.URL, resp.StatusCode)
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func checkMedusaHealth(store models.Store) (bool, error) {
	// For now, just check if the service is accessible
	// In a real implementation, you'd check the Medusa API health endpoint
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	healthURL := store.URL + "/health"
	resp, err := client.Get(healthURL)
	if err != nil {
		log.Printf("Health check failed for %s: %v", healthURL, err)
		return false, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		log.Printf("Health check passed for %s", healthURL)
		return true, nil
	}
	
	log.Printf("Health check failed for %s - status: %d", healthURL, resp.StatusCode)
	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// GetStorePodStatus gets the actual pod status from Kubernetes
func GetStorePodStatus(store models.Store) (string, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	
	cmd := exec.Command("kubectl", "get", "pods", 
		"--namespace", store.Namespace,
		"--selector", "app.kubernetes.io/name=wordpress", // For WooCommerce
		"--output", "jsonpath={.items[*].status.phase}",
		"--kubeconfig", kubeconfig)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to get pod status for %s: %v", store.Namespace, err)
		return "", err
	}
	
	phases := strings.Fields(strings.TrimSpace(string(output)))
	if len(phases) == 0 {
		return "Unknown", nil
	}
	
	// Return the most common phase or the first non-failed one
	for _, phase := range phases {
		if phase == "Running" {
			return "Running", nil
		}
		if phase == "Pending" {
			return "Pending", nil
		}
	}
	
	return phases[0], nil
}

// WaitForStoreReady waits for a store to become ready with timeout
func WaitForStoreReady(store models.Store, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	log.Printf("Waiting for store %s to become ready...", store.ID)
	
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for store %s to become ready", store.ID)
		case <-ticker.C:
			healthy, err := CheckStoreHealth(store)
			if err != nil {
				log.Printf("Health check failed for store %s: %v", store.ID, err)
				continue
			}
			
			if healthy {
				log.Printf("Store %s is now ready and healthy", store.ID)
				return nil
			}
		}
	}
}

// ReconcileStoreStatus checks the actual state and updates the database if needed
func ReconcileStoreStatus(store models.Store, db interface{}) error {
	// Get actual pod status
	podStatus, err := GetStorePodStatus(store)
	if err != nil {
		log.Printf("Failed to get pod status for store %s: %v", store.ID, err)
		return err
	}
	
	// Determine what the status should be
	expectedStatus := map[string]string{
		"Running":  "Ready",
		"Pending":  "Provisioning",
		"Failed":   "Failed",
		"Unknown":  "Provisioning",
	}[podStatus]
	
	// If the status is different, update it
	if expectedStatus != store.Status {
		log.Printf("Reconciling store %s status: %s -> %s", store.ID, store.Status, expectedStatus)
		
		// This would need to be implemented based on your DB interface
		// For now, just log the change
		// db.Model(&store).Update("status", expectedStatus)
	}
	
	return nil
}
