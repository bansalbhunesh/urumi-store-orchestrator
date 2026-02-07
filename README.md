# 🛍️ Urumi Store Orchestrator (Round 1)

**A production-ready Kubernetes platform for provisioning e-commerce stores on demand.**

This project implements a "Shopify-like" provisioning system where users can create isolated WooCommerce stores instantly. It runs on local Kubernetes (Kind) but is designed to ship to Production (VPS/k3s) using the exact same Helm charts.

![Dashboard Preview](dashboard/preview.png)

## ✨ Key Features

- **🚀 Instant Provisioning**: Launches WooCommerce (Full Stack) or Medusa (Simulated/Lightweight) stores.
- **🔒 Strong Isolation**: Each store runs in its own **Kubernetes Namespace** (`store-<uuid>`).
- **🌐 Automatic Ingress**: Assigns unique URLs (e.g., `http://store-abc.localhost`) automatically.
- **📦 Helm-Native**: Uses standard Helm charts for deployment, ensuring portability between Local (Kind) and Production (k3s).
- **🎨 Modern Dashboard**: Beautiful React UI with real-time status polling.

---

## 🛠️ Prerequisites

Before you start, ensure you have:
1. **Docker Desktop** (Running) - *Required for the local cluster.*
2. **Go** (v1.21+) - *For the backend.*
3. **Node.js** (v18+) - *For the dashboard.*
4. **Kind** & **Kubectl** - *(Optional, the app can install Kind for you).*

---

## 🚀 Quick Start (Local)

### 1. Start the Kubernetes Cluster
If you don't have a cluster yet, create one with Ingress support:
```bash
# From the project root
kind create cluster --config kind-config.yaml --name urumi-cluster
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

### 2. Start the Backend (Orchestrator)
The backend listens on `:8080`.
```bash
cd backend
go mod tidy
# Run the server
go run .
# Or build and run
# go build -o orchestrator.exe && .\orchestrator.exe
```

### 3. Start the Dashboard
The frontend runs on `:5173`.
```bash
cd dashboard
npm install
npm run dev
```

### 4. Provision a Store!
1. Open [http://localhost:5173](http://localhost:5173).
2. Click **New Store**.
3. Select **WooCommerce** and give it a name.
4. Watch the status go from `Provisioning` ➡️ `Ready`.
5. Click **Visit Store** to see your live WooCommerce site!

---

## 🏗️ Architecture & System Design

### The "Control Plane"
- **Backend (Go)**: Acts as the intelligent orchestrator. It wraps the **Helm SDK** to execute deployments programmatically. It maintains state in a lightweight embedded SQLite DB.
- **Frontend (React)**: Polling-based UI that queries the Go API for store status.

### The "Data Plane" (Kubernetes)
- **Namespace Isolation**: Every store gets a dedicated namespace. This ensures that even if one store is compromised or crashes, others are unaffected.
- **Persistence**: We use `PersistentVolumeClaims` (PVCs) for both the MariaDB database and the WordPress file system.
- **Ingress Layer**: Nginx Ingress Controller routes traffic based on hostnames (`*.localhost`).

### Local vs. Production Strategy
We solve the "Local to Prod" challenge using Helm Value overlays:

| Feature | Local (`values-local.yaml`) | Production (`values-prod.yaml`) |
| :--- | :--- | :--- |
| **Domain** | `*.localhost` | `*.urumi.ai` (or similar) |
| **Storage** | Standard (HostPath) | `local-path` / Longhorn |
| **TLS** | Disabled (HTTP) | Enabled (Cert-Manager) |

---

## 🧪 Verification
You can verify the system internals using `kubectl`:

```bash
# List all running stores (namespaces)
kubectl get ns

# Check pods for a specific store
kubectl get pods -n store-<store-id>

# Check generated Ingress rules
kubectl get ingress -A
```

---

## 📁 Project Structure

```
urumi-store-orchestrator/
├── backend/          # Go Orchestrator API
├── dashboard/        # React + Vite Frontend
├── charts/           # Helm Charts
│   └── woocommerce/  # The Store Blueprint
└── kind-config.yaml  # Cluster setup config
```

---

*Verified for Urumi AI SDE Internship - Round 1*
