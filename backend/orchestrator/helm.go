package orchestrator

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
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

	// Generate Random Passwords
	rootPass := generatePassword(16)
	dbPass := generatePassword(16)
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

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error provisioning store %s: %s\nOutput: %s\n", store.ID, err, string(output))
		return err
	}

	fmt.Printf("Successfully provisioned store %s at %s\n", store.ID, host)
	return nil
}

func generatePassword(length int) string {
	// fast and dirty random string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// DeleteStore runs the helm uninstall command
func DeleteStore(store models.Store) error {
	// Get KUBECONFIG
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	cmd := exec.Command("helm", "uninstall", store.Namespace, "--namespace", store.Namespace, "--kubeconfig", kubeconfig)

	// Also delete namespace? helm uninstall doesn't delete namespace usually.
	// Let's delete the namespace directly.

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error uninstalling helm release %s: %s\nOutput: %s\n", store.ID, err, string(output))
		// continue to try deleting namespace
	}

	cmdNs := exec.Command("kubectl", "delete", "namespace", store.Namespace, "--kubeconfig", kubeconfig)
	outputNs, errNs := cmdNs.CombinedOutput()
	if errNs != nil {
		fmt.Printf("Error deleting namespace %s: %s\nOutput: %s\n", store.Namespace, errNs, string(outputNs))
		return errNs
	}

	return nil
}
