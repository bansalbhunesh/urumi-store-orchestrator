Write-Host "Starting Urumi Store Orchestrator..." -ForegroundColor Green

# 1. Generate Kubeconfig
Write-Host "Exporting Kubeconfig..."
try {
    if (Test-Path ".\kubeconfig.yaml") {
        Remove-Item ".\kubeconfig.yaml" -Force
    }
    kind export kubeconfig --name urumi-cluster --kubeconfig .\kubeconfig.yaml
    $env:KUBECONFIG = "$PWD\kubeconfig.yaml"
} catch {
    Write-Host "Warning: Could not export kubeconfig. Ensure Kind is running." -ForegroundColor Yellow
}

# 1.5 Ensure Ingress Controller
Write-Host "Checking Ingress Controller..."
$ingress = kubectl get namespace ingress-nginx --ignore-not-found
if (-not $ingress) {
    Write-Host "Installing Ingress NGINX..."
    if (-not (Test-Path "ingress-nginx.yaml")) {
        Invoke-WebRequest -Uri https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml -OutFile ingress-nginx.yaml
    }
    kubectl apply -f ingress-nginx.yaml
    Write-Host "Ingress Controller installed."
} else {
    Write-Host "Ingress Controller already exists."
}

# 2. Start Backend
Write-Host "Launching Backend (Port 8080)..."
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd backend; `$env:KUBECONFIG='$PWD\kubeconfig.yaml'; go run ."

# 2. Start Frontend
Write-Host "Launching Dashboard (Port 5173)..."
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd dashboard; npm run dev"

# 3. Open Browser
Start-Sleep -Seconds 5
Start-Process "http://localhost:5173"

Write-Host "Systems operational!"
Write-Host "Backend: http://localhost:8080"
Write-Host "Frontend: http://localhost:5173"
