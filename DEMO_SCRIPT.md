# Demo Video Script 🎬

## 1. Intro (30s)
- **Show**: This `README.md` on GitHub.
- **Say**: "Hi, I'm Bhunesh. This is my submission for the Urumi AI SDE Internship. I've built a Kubernetes-native Store Provisioning Platform that supports both WooCommerce and Medusa."

## 2. Infrastructure Walkthrough (1m)
- **Show**: VS Code -> `charts/woocommerce` and `charts/medusa`.
- **Say**: "I use Helm charts for defining the store infrastructure. This ensures my deployment logic is identical for both Local Kind clusters and Production."
- **Show**: `backend/orchestrator/helm.go`.
- **Say**: "My Go backend acts as an orchestrator, dynamically selecting the correct Helm chart based on user input."

## 3. The "Wow" Factor - Provisioning (1m)
- **Action**: Open Dashboard at `http://localhost:5173`.
- **Action**: Click **New Store**. Select **WooCommerce**. Name: "My Shop".
- **Action**: While it spins, switch to Terminal.
- **Command**: `kubectl get pods -n store-<id> -w`
- **Say**: "You can see the backend has triggered a real Kubernetes deployment in its own isolated Namespace."

## 4. Medusa Support (30s)
- **Action**: Create another store. Select **Medusa**.
- **Say**: "The platform is extensible. Here I'm provisioning a Medusa store. For this demo, I'm using a lightweight version to save resources on my local machine, but the orchestration flow is exactly the same."

## 5. Cleanup (30s)
- **Action**: Click **Delete** on the dashboard.
- **Command**: `kubectl get ns` -> Show the namespace failing.
- **Say**: "Thank you!"
