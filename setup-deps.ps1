# Setup Helm Dependencies for Urumi Store Orchestrator
Write-Host "🔧 Setting up Helm dependencies..." -ForegroundColor Yellow

# Add required Helm repositories
Write-Host "📦 Adding Bitnami repository..." -ForegroundColor Cyan
helm repo add bitnami https://charts.bitnami.com/bitnami
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to add Bitnami repository" -ForegroundColor Red
    exit 1
}

# Update repositories
Write-Host "🔄 Updating Helm repositories..." -ForegroundColor Cyan
helm repo update
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to update repositories" -ForegroundColor Red
    exit 1
}

# Build dependencies for WooCommerce chart
Write-Host "🛒 Building WooCommerce chart dependencies..." -ForegroundColor Cyan
Set-Location charts\woocommerce
helm dependency build
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to build WooCommerce dependencies" -ForegroundColor Red
    Set-Location ..\..
    exit 1
}
Set-Location ..\..

# Build dependencies for Medusa chart
Write-Host "🚀 Building Medusa chart dependencies..." -ForegroundColor Cyan
Set-Location charts\medusa
helm dependency build
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to build Medusa dependencies" -ForegroundColor Red
    Set-Location ..\..
    exit 1
}
Set-Location ..\..

Write-Host "✅ All Helm dependencies setup successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "🎉 You can now create stores from the dashboard!" -ForegroundColor Green
Write-Host "🌐 Dashboard: http://localhost:5173" -ForegroundColor Cyan
