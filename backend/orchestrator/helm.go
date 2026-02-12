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

	// Assuming running from backend/ directory
	// Determine chart based on store type
	var chartPath string
	// Default values file, can be overridden by env
	defaultValuesFile := "../charts/woocommerce/values-local.yaml"

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	domainSuffix := os.Getenv("DOMAIN_SUFFIX")
	if domainSuffix == "" {
		domainSuffix = "localhost"
	}

	valuesFile := os.Getenv("HELM_VALUES_FILE")

	switch store.Type {
	case "medusa":
		chartPath = "../charts/medusa"
		if valuesFile == "" {
			valuesFile = "../charts/medusa/values-local.yaml"
		}
	default:
		chartPath = "../charts/woocommerce"
		if valuesFile == "" {
			valuesFile = defaultValuesFile
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

	host := fmt.Sprintf("%s.%s", store.Namespace, domainSuffix)

	cmd := exec.Command("helm", "upgrade", "--install", releaseName, chartPath,
		"--kubeconfig", kubeconfig,
		"--namespace", store.Namespace,
		"--create-namespace",
		"--values", valuesFile,
		"--set", fmt.Sprintf("ingress.hosts[0].host=%s", host),
		"--set", fmt.Sprintf("mariadb.auth.rootPassword=%s", rootPass),
		"--set", fmt.Sprintf("mariadb.auth.password=%s", dbPass),
		"--set", fmt.Sprintf("wordpress.db.password=%s", dbPass), // DB Connection pass
		"--set", fmt.Sprintf("wordpress.password=%s", wpInternalPass), // Admin Panel pass
		// woocommerce values.yaml:
		// mariadb.auth.password
		// mariadb.auth.rootPassword
		// wordpress.password is NOT there, it uses env WORDPRESS_DB_PASSWORD from values.mariadb.auth.password
		// But wait, the deployment.yaml uses:
		// value: {{ .Values.mariadb.auth.password | quote }}
		// So setting 'mariadb.auth.password' is enough for both DB and WordPress connection?
		// Let's verify deployment.yaml env vars.
		// WORDPRESS_DB_PASSWORD value: {{ .Values.mariadb.auth.password | quote }}
		// Yes.
	)

	// We need to pass the DB password to the MariaDB chart AND the WordPress deployment.
	// typically the wordpress chart depends on mariadb.
	// If we set mariadb.auth.password, it updates the MariaDB chart.
	// The WordPress deployment reads from .Values.mariadb.auth.password.
	// So just setting mariadb.auth.password should work for both.

	log.Printf("Executing helm command for store %s: %s", store.ID, cmd.String())
	
	// Set timeout and prevent hanging
	cmd.Start()
	
	// Create a channel to receive completion signal
	done := make(chan error, 1)
	
	// Wait for command to complete in goroutine
	go func() {
		done <- cmd.Wait()
	}()
	
	// Wait for completion or timeout
	select {
	case err := <-done:
		output, _ := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error provisioning store %s: %s\nOutput: %s\n", store.ID, err, string(output))
			return fmt.Errorf("helm install failed: %w", err)
		}
		log.Printf("Successfully provisioned store %s at %s\nOutput: %s", store.ID, host, string(output))
		return nil
	case <-time.After(5 * time.Minute):
		// Kill the process if it hangs
		cmd.Process.Kill()
		return fmt.Errorf("helm install timed out after 5 minutes for store %s", store.ID)
	}
}

func generateSecurePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
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
