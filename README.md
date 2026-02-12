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
- **🛡️ Security First**: Rate limiting, CORS protection, input validation, and security headers.
- **📊 Real-time Monitoring**: Health checks, status reconciliation, and comprehensive logging.
- **🔄 Error Recovery**: Automatic retry mechanisms and detailed error reporting.

---

## 🛠️ Prerequisites

Before you start, ensure you have:
1. **Docker Desktop** (Running) - *Required for the local cluster.*
2. **Go** (v1.21+) - *For the backend.*
3. **Node.js** (v18+) - *For the dashboard.*
4. **Kind** & **Kubectl** - *(Optional, the app can install Kind for you).*
5. **Helm** (v3+) - *For chart deployment.*

---

## 🚀 Quick Start (Local)

### 1. Start the Kubernetes Cluster
If you don't have a cluster yet, create one with Ingress support:
```bash
# From the project root
kind create cluster --config kind-config.yaml --name urumi-cluster
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

### 2. One-Click Start (Windows)
Run the helper script to launch both Backend and Frontend:
```powershell
.\test-e2e.ps1
```

### 3. Manual Start
#### Backend
```bash
cd backend
go mod tidy
go build -o urumi-backend.exe .
.\urumi-backend.exe
```

#### Frontend
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

## 🧪 Testing

Run the comprehensive end-to-end test suite:
```powershell
.\test-e2e.ps1
```

This script tests:
- ✅ Prerequisites and cluster setup
- ✅ Backend API endpoints
- ✅ Store creation and deletion
- ✅ Rate limiting and security features
- ✅ CORS and health checks
- ✅ Frontend build process

---

## 🏗️ Architecture & System Design

### The "Control Plane"
- **Backend (Go)**: Acts as the intelligent orchestrator with:
  - RESTful API with comprehensive error handling
  - Background reconciliation service for status sync
  - Rate limiting and security middleware
  - Health checks and monitoring
  - Secure password generation
- **Frontend (React)**: Modern UI with:
  - Real-time status updates
  - Error boundaries and graceful error handling
  - Loading states and user feedback
  - Responsive design with Tailwind CSS

### The "Data Plane" (Kubernetes)
- **Namespace Isolation**: Every store gets a dedicated namespace
- **Persistence**: PVCs for MariaDB database and WordPress files
- **Ingress Layer**: Nginx Ingress Controller with hostname-based routing
- **Health Monitoring**: Liveness and readiness probes
- **Resource Management**: CPU/memory limits and requests

### Security Features
- **Rate Limiting**: Token bucket algorithm (10 req/min, burst 20)
- **CORS Protection**: Configurable origin allowlist
- **Input Validation**: Sanitization and length limits
- **Security Headers**: CSP, XSS protection, frame options
- **Secure Passwords**: Cryptographically random generation

### Local vs. Production Strategy
We solve the "Local to Prod" challenge using Helm Value overlays:

| Feature | Local (`values-local.yaml`) | Production (`values-prod.yaml`) |
| :--- | :--- | :--- |
| **Domain** | `*.localhost` | `*.example.com` |
| **Storage** | Standard (HostPath) | `local-path` / Longhorn |
| **TLS** | Disabled (HTTP) | Enabled (Cert-Manager) |
| **Resources** | Development limits | Production quotas |
| **Image Tags** | `latest` | Pinned versions |

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

# Check store health
curl http://localhost:8080/api/stores/<store-id>/health

# View backend logs
curl http://localhost:8080/health
```

---

## 📁 Project Structure

```
urumi-store-orchestrator/
├── backend/                 # Go Orchestrator API
│   ├── handlers/           # API request handlers
│   ├── middleware/         # Security & rate limiting
│   ├── models/            # Data models
│   ├── orchestrator/      # Helm & K8s operations
│   └── main.go           # Application entry point
├── dashboard/              # React + Vite Frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   └── App.jsx      # Main application
│   └── package.json
├── charts/                 # Helm Charts
│   ├── woocommerce/      # Store blueprint
│   └── medusa/          # Medusa stub (Round 2)
├── test-e2e.ps1          # End-to-end test suite
└── kind-config.yaml       # Cluster setup config
```

---

## 🔧 Configuration

### Environment Variables
- `DOMAIN_SUFFIX`: Override default domain suffix (default: `localhost`)
- `HELM_VALUES_FILE`: Override default Helm values file
- `KUBECONFIG`: Path to kubeconfig file
- `ALLOWED_ORIGINS`: Comma-separated list of allowed CORS origins

### Security Configuration
- Rate limiting: 10 requests/minute per IP (burst: 20)
- Request timeout: 30 seconds
- Max request size: 1MB
- CORS: Strict origin validation

---

## 🚀 Production Deployment

### VPS/k3s Setup
1. Install k3s: `curl -sfL https://get.k3s.io | sh -`
2. Install Helm: Follow official Helm installation guide
3. Configure domain DNS to point to your VPS
4. Update `values-prod.yaml` with your domain
5. Deploy: `helm install urumi ./charts/urumi -f values-prod.yaml`

### Monitoring & Observability
- **Health Checks**: `/health` endpoint and per-store health monitoring
- **Logging**: Structured logging with request tracking
- **Metrics**: Built-in reconciliation status and error tracking
- **Status Reconciliation**: Background service ensures state consistency

---

## 🛡️ Security Considerations

- **Input Validation**: All user inputs are validated and sanitized
- **Rate Limiting**: Prevents abuse and resource exhaustion
- **Namespace Isolation**: Complete tenant separation
- **Secrets Management**: No hardcoded secrets, secure password generation
- **Network Policies**: Ready for implementation (chart supports)
- **RBAC**: Principle of least privilege (can be extended)

---

## 🔄 Troubleshooting

### Common Issues
1. **Store stuck in Provisioning**: Check pod logs with `kubectl logs -n store-<id>`
2. **Ingress not working**: Verify ingress controller is running
3. **Database connection failed**: Check PVC status and MariaDB logs
4. **Rate limiting errors**: Wait for tokens to replenish (1 per 6 seconds)

### Debug Commands
```bash
# Check backend logs
curl http://localhost:8080/health

# Check specific store
kubectl get all -n store-<id>

# View Helm release
helm list -n store-<id>

# Check ingress rules
kubectl describe ingress -n store-<id>
```

---

*Verified for Urumi AI SDE Internship - Round 1*
