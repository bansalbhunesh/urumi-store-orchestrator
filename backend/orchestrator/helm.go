package orchestrator

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"urumi-backend/models"
)

// ProvisionStore runs the helm install command
func ProvisionStore(store models.Store) error {
	// Helm install command
	// helm install <release-name> ../charts/woocommerce --namespace <ns> --create-namespace --set ...

	// Determine chart based on store type
	var chartPath string
	
	// Resolve chart directory relative to current working directory or executable
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	
	// Check if running from root or backend
	if _, err := os.Stat(filepath.Join(baseDir, "charts")); err == nil {
		// Running from root
		chartPath = filepath.Join(baseDir, "charts")
	} else if _, err := os.Stat(filepath.Join(baseDir, "../charts")); err == nil {
		// Running from backend
		chartPath = filepath.Join(baseDir, "../charts")
	} else {
		// Fallback to Env var or assume relative
		chartPath = os.Getenv("CHARTS_DIR")
		if chartPath == "" {
			return fmt.Errorf("could not locate charts directory")
		}
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Domain Suffix Logic
	domainSuffix := os.Getenv("DOMAIN_SUFFIX")
	if domainSuffix == "" {
		domainSuffix = "localhost" // Default for local dev
	}

	// Values File Logic
	valuesFile := os.Getenv("HELM_VALUES_FILE")
	
	var specificChartPath string

	switch store.Type {
	case "medusa":
		specificChartPath = filepath.Join(chartPath, "medusa")
		if valuesFile == "" {
			valuesFile = filepath.Join(specificChartPath, "values-local.yaml")
		}
	default:
		specificChartPath = filepath.Join(chartPath, "woocommerce")
		if valuesFile == "" {
			valuesFile = filepath.Join(specificChartPath, "values-local.yaml")
		}
	}

	releaseName := store.Namespace

	// Generate secure random passwords
	rootPass, err := generateSecurePassword(16)
	if err != nil {
		return fmt.Errorf("failed to generate root password: %w", err)
	}
	dbPass, err := generateSecurePassword(16)
	if err != nil {
		return fmt.Errorf("failed to generate database password: %w", err)
	}
	// For demo purposes, we set a known admin password so the user can login
	wpInternalPass := "password123"

	// Host generation: store-uuid.domain or just store-uuid for simple setups
	// If domainSuffix is "localhost", we might want store-uuid.localhost
	host := fmt.Sprintf("%s.%s", store.Namespace, domainSuffix)

	cmd := exec.Command("helm", "upgrade", "--install", releaseName, specificChartPath,
		"--kubeconfig", kubeconfig,
		"--namespace", store.Namespace,
		"--create-namespace",
		"--values", valuesFile,
		"--set", fmt.Sprintf("ingress.hosts[0].host=%s", host),
		"--set", fmt.Sprintf("mariadb.auth.rootPassword=%s", rootPass),
		"--set", fmt.Sprintf("mariadb.auth.password=%s", dbPass),
		"--set", fmt.Sprintf("wordpress.db.password=%s", dbPass), // DB Connection pass
		"--set", fmt.Sprintf("wordpress.password=%s", wpInternalPass), // Admin Panel pass
		"--wait", // Wait for resources to be ready (optional, might timeout long operations)
		"--timeout", "10m",
	)

	log.Printf("Executing helm command for store %s: %s", store.ID, cmd.String())
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error provision store %s: %v\nOutput: %s", store.ID, err, string(output))
		return fmt.Errorf("helm install failed: %w - Output: %s", err, string(output))
	}
	
	log.Printf("Successfully provisioned store %s at %s\nOutput: %s", store.ID, host, string(output))
	return nil
}

func generateSecurePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Removed special chars to avoid shell escaping issues
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

// DeleteStore runs the helm uninstall command
func DeleteStore(store models.Store) error {
	// Get KUBECONFIG
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	log.Printf("Starting deletion of store %s (%s)", store.ID, store.Name)

	// First, try to uninstall the helm release
	cmd := exec.Command("helm", "uninstall", store.Namespace, "--namespace", store.Namespace, "--kubeconfig", kubeconfig)
	log.Printf("Executing helm uninstall for store %s: %s", store.ID, cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error uninstalling helm release %s: %s\nOutput: %s\n", store.ID, err, string(output))
		// Continue to try deleting namespace anyway
	} else {
		log.Printf("Successfully uninstalled helm release for store %s", store.ID)
	}

	// Wait a moment for resources to be cleaned up
	time.Sleep(2 * time.Second)

	// Delete the namespace
	cmdNs := exec.Command("kubectl", "delete", "namespace", store.Namespace, "--kubeconfig", kubeconfig)
	log.Printf("Executing kubectl delete namespace for store %s: %s", store.ID, cmdNs.String())
	outputNs, errNs := cmdNs.CombinedOutput()
	if errNs != nil {
		log.Printf("Error deleting namespace %s: %s\nOutput: %s\n", store.Namespace, errNs, string(outputNs))
		return fmt.Errorf("failed to delete namespace: %w", errNs)
	}

	log.Printf("Successfully deleted namespace %s for store %s", store.Namespace, store.ID)
	return nil
}
