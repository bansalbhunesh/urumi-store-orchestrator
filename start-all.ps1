# Start setup
Write-Host "🚀 Starting Urumi Store Orchestrator..." -ForegroundColor Green

# 1. Start Backend
Write-Host "📦 Launching Backend (Port 8080)..."
Start-Process -FilePath "powershell" -ArgumentList "-NoExit", "-Command", "& {cd backend; go run .}" -WorkingDirectory "."

# 2. Start Frontend
Write-Host "🎨 Launching Dashboard (Port 5173)..."
Start-Process -FilePath "powershell" -ArgumentList "-NoExit", "-Command", "& {cd dashboard; npm run dev}" -WorkingDirectory "."

# 3. Open Browser
Start-Sleep -Seconds 5
Start-Process "http://localhost:5173"

Write-Host "✅ Systems operational!"
Write-Host "Backend: http://localhost:8080"
Write-Host "Frontend: http://localhost:5173"
