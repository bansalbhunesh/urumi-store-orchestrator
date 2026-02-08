package orchestrator

import (
	"fmt"
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
	var valuesFile string

	switch store.Type {
	case "medusa":
		chartPath = "../charts/medusa"
		valuesFile = "../charts/medusa/values-local.yaml"
	default:
		chartPath = "../charts/woocommerce"
		valuesFile = "../charts/woocommerce/values-local.yaml"
	}

	releaseName := store.Namespace

	// Get KUBECONFIG from env or use default
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		// Fallback to default location if not set, but better to be explicit in our setup
		home, _ := os.UserHomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	cmd := exec.Command("helm", "upgrade", "--install", releaseName, chartPath,
		"--kubeconfig", kubeconfig,
		"--namespace", store.Namespace,
		"--create-namespace",
		"--values", valuesFile,
		"--set", fmt.Sprintf("ingress.hosts[0].host=%s.localhost", store.Namespace),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error provisioning store %s: %s\nOutput: %s\n", store.ID, err, string(output))
		return err
	}

	fmt.Printf("Successfully provisioned store %s\n", store.ID)
	return nil
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
