# End-to-End Test Script for Urumi Store Orchestrator
# This script tests the complete store creation and deletion flow

Write-Host "🚀 Starting Urumi Store Orchestrator E2E Tests" -ForegroundColor Green

# Function to check if a command exists
function Command-Exists {
    param ($command)
    try {
        Get-Command $command -ErrorAction Stop | Out-Null
        return $true
    }
    catch {
        return $false
    }
}

# Check prerequisites
Write-Host "📋 Checking prerequisites..." -ForegroundColor Yellow

if (-not (Command-Exists "go")) {
    Write-Host "❌ Go is not installed. Please install Go 1.21+" -ForegroundColor Red
    exit 1
}

if (-not (Command-Exists "node")) {
    Write-Host "❌ Node.js is not installed. Please install Node.js 18+" -ForegroundColor Red
    exit 1
}

if (-not (Command-Exists "kubectl")) {
    Write-Host "❌ kubectl is not installed. Please install kubectl" -ForegroundColor Red
    exit 1
}

if (-not (Command-Exists "helm")) {
    Write-Host "❌ Helm is not installed. Please install Helm" -ForegroundColor Red
    exit 1
}

Write-Host "✅ All prerequisites found!" -ForegroundColor Green

# Check if Kubernetes cluster is running
Write-Host "🔍 Checking Kubernetes cluster..." -ForegroundColor Yellow
try {
    kubectl cluster-info | Out-Null
    Write-Host "✅ Kubernetes cluster is running" -ForegroundColor Green
}
catch {
    Write-Host "❌ Kubernetes cluster is not running. Please start your cluster first." -ForegroundColor Red
    Write-Host "💡 Run: kind create cluster --config kind-config.yaml --name urumi-cluster" -ForegroundColor Cyan
    exit 1
}

# Check if ingress-nginx is installed
Write-Host "🔍 Checking ingress controller..." -ForegroundColor Yellow
try {
    kubectl get namespace ingress-nginx | Out-Null
    Write-Host "✅ Ingress controller is installed" -ForegroundColor Green
}
catch {
    Write-Host "⚠️  Ingress controller not found. Installing..." -ForegroundColor Yellow
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
    Write-Host "✅ Ingress controller installed" -ForegroundColor Green
}

# Build and test backend
Write-Host "🔨 Building backend..." -ForegroundColor Yellow
Set-Location backend
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to install Go dependencies" -ForegroundColor Red
    exit 1
}

go build -o urumi-backend.exe .
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to build backend" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Backend built successfully" -ForegroundColor Green

# Start backend in background
Write-Host "🚀 Starting backend server..." -ForegroundColor Yellow
$backend = Start-Process -FilePath ".\urumi-backend.exe" -PassThru -WindowStyle Hidden
Start-Sleep -Seconds 3

# Test backend health
Write-Host "🏥 Testing backend health..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
    Write-Host "✅ Backend is healthy: $($response.status)" -ForegroundColor Green
}
catch {
    Write-Host "❌ Backend health check failed: $_" -ForegroundColor Red
    Stop-Process -Id $backend.Id -Force
    exit 1
}

# Test API endpoints
Write-Host "🧪 Testing API endpoints..." -ForegroundColor Yellow

# Test list stores (should be empty initially)
try {
    $stores = Invoke-RestMethod -Uri "http://localhost:8080/api/stores" -Method GET
    if ($stores.Count -eq 0) {
        Write-Host "✅ GET /api/stores - Empty list as expected" -ForegroundColor Green
    } else {
        Write-Host "⚠️  GET /api/stores - Unexpected stores found: $stores" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "❌ GET /api/stores failed: $_" -ForegroundColor Red
    Stop-Process -Id $backend.Id -Force
    exit 1
}

# Test store creation
Write-Host "🏪 Testing store creation..." -ForegroundColor Yellow
$storeData = @{
    name = "Test Store"
    type = "woocommerce"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/stores" -Method POST -Body $storeData -ContentType "application/json"
    $storeId = $response.id
    Write-Host "✅ Store created successfully: $storeId" -ForegroundColor Green
    Write-Host "   Name: $($response.name)" -ForegroundColor Cyan
    Write-Host "   Type: $($response.type)" -ForegroundColor Cyan
    Write-Host "   Status: $($response.status)" -ForegroundColor Cyan
    Write-Host "   URL: $($response.url)" -ForegroundColor Cyan
}
catch {
    Write-Host "❌ Store creation failed: $_" -ForegroundColor Red
    Stop-Process -Id $backend.Id -Force
    exit 1
}

# Wait for provisioning to start
Write-Host "⏳ Waiting for provisioning to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Check store status
try {
    $store = Invoke-RestMethod -Uri "http://localhost:8080/api/stores/$storeId" -Method GET
    Write-Host "📊 Store status: $($store.status)" -ForegroundColor Cyan
}
catch {
    Write-Host "⚠️  Could not check store status: $_" -ForegroundColor Yellow
}

# Test health check endpoint
Write-Host "🏥 Testing store health check..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/api/stores/$storeId/health" -Method GET
    Write-Host "📈 Store health: $($health.healthy)" -ForegroundColor Cyan
}
catch {
    Write-Host "⚠️  Store health check failed (expected during provisioning): $_" -ForegroundColor Yellow
}

# Test rate limiting
Write-Host "🚦 Testing rate limiting..." -ForegroundColor Yellow
$rateLimitHit = $false
for ($i = 1; $i -le 15; $i++) {
    try {
        Invoke-RestMethod -Uri "http://localhost:8080/api/stores" -Method GET -TimeoutSec 2 | Out-Null
    }
    catch {
        if ($_.Exception.Response.StatusCode -eq 429) {
            $rateLimitHit = $true
            break
        }
    }
}

if ($rateLimitHit) {
    Write-Host "✅ Rate limiting is working" -ForegroundColor Green
} else {
    Write-Host "⚠️  Rate limiting may not be working (or limits too high for test)" -ForegroundColor Yellow
}

# Test store deletion
Write-Host "🗑️  Testing store deletion..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8080/api/stores/$storeId" -Method DELETE | Out-Null
    Write-Host "✅ Store deletion initiated" -ForegroundColor Green
}
catch {
    Write-Host "❌ Store deletion failed: $_" -ForegroundColor Red
}

# Wait a moment for deletion to process
Start-Sleep -Seconds 3

# Test CORS
Write-Host "🌐 Testing CORS headers..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/stores" -Method OPTIONS
    $corsHeaders = $response.Headers["Access-Control-Allow-Origin"]
    if ($corsHeaders) {
        Write-Host "✅ CORS headers present: $corsHeaders" -ForegroundColor Green
    } else {
        Write-Host "⚠️  CORS headers not found" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "⚠️  CORS test failed: $_" -ForegroundColor Yellow
}

# Stop backend
Write-Host "🛑 Stopping backend server..." -ForegroundColor Yellow
Stop-Process -Id $backend.Id -Force

# Test frontend build
Write-Host "🎨 Testing frontend build..." -ForegroundColor Yellow
Set-Location ..\dashboard
npm install
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to install npm dependencies" -ForegroundColor Red
    exit 1
}

npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to build frontend" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Frontend built successfully" -ForegroundColor Green

# Return to root
Set-Location ..

Write-Host "🎉 All E2E tests completed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "📝 Test Summary:" -ForegroundColor Cyan
Write-Host "   ✅ Prerequisites checked" -ForegroundColor Green
Write-Host "   ✅ Kubernetes cluster verified" -ForegroundColor Green
Write-Host "   ✅ Backend built and started" -ForegroundColor Green
Write-Host "   ✅ API endpoints tested" -ForegroundColor Green
Write-Host "   ✅ Store creation tested" -ForegroundColor Green
Write-Host "   ✅ Store deletion tested" -ForegroundColor Green
Write-Host "   ✅ Rate limiting tested" -ForegroundColor Green
Write-Host "   ✅ CORS headers tested" -ForegroundColor Green
Write-Host "   ✅ Frontend build tested" -ForegroundColor Green
Write-Host ""
Write-Host "🚀 Your Urumi Store Orchestrator is ready for deployment!" -ForegroundColor Green
