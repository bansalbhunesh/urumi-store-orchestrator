# System Design & Tradeoffs

## Architecture Choice
This system follows a **Control Plane** architecture where a centralized Orchestrator (Go Backend) manages the lifecycle of tenant resources (Stores) on the Kubernetes cluster.

### Components
1. **Dashboard (React)**: User-facing UI for managing stores. Polls the API for state changes.
2. **Orchestrator (Go)**: 
   - Exposes REST API.
   - Maintains state in a lightweight SQLite database (for fast retrieval without querying K8s API constantly).
   - Wraps `helm` binary to interface with the Kubernetes API Server.
3. **Infrastructure (Helm)**:
   - **WooCommerce Chart**: Bundles WordPress and MariaDB.
   - **Isolation**: Each store is deployed into a dedicated Namespace (e.g., `store-abc1234`).

## Local vs Production Strategy
We use a "Single Chart, Multiple Values" strategy.
- **Local (Kind)**: `values-local.yaml` optimizes for speed and local accessibility (NodePort/Localhost Ingress, Standard StorageClass).
- **Production (VPS)**: `values-prod.yaml` targets k3s defaults (Local-Path StorageClass, Traefik Ingress) and enables features like TLS (via cert-manager, conceptually).

## Tradeoffs

### 1. Helm Wrapper vs Kubernetes Operator (CRDs)
- **Decision**: We chose the **Helm Wrapper** approach.
- **Why**: It is significantly faster to implement and easier to debug for "Day 1". It allows using standard Helm tools to inspect releases.
- **Tradeoff**: We lose the continuous reconciliation loop that a K8s Operator provides. If a store pod is deleted manually, the Orchestrator doesn't know until we query.
- **Mitigation**: A background "Sync/Reconciliation" loop could be added to the Go backend to verify actual K8s state matches DB state.

### 2. Database (SQLite)
- **Decision**: Embedded SQLite.
- **Why**: Zero-dep startup for the evaluator. No need to provision a Postgres RDS just for the control plane.
- **Tradeoff**: Not horizontally scalable.
- **Mitigation**: Use Postgres in production (change GORM driver).

### 3. Namespace per Store
- **Decision**: Hard isolation.
- **Why**: Best security and resource management (easier to delete, quota management).
- **Tradeoff**: Higher overhead on the cluster (more namespaces).

## Future Improvements
- **Rate Limiting**: Implement token bucket in the Go API.
- **Auth**: Add JWT authentication for separating user stores.
- **Observability**: Prometheus metrics for provisioning duration and failure rates.
