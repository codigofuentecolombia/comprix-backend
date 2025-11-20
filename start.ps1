# Script para iniciar el servidor Comprix
Write-Host "==================================" -ForegroundColor Cyan
Write-Host "  Iniciando Servidor Comprix" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Agregar Go al PATH
$env:PATH += ";C:\Program Files\Go\bin"

# Cambiar al directorio del proyecto
Set-Location -Path $PSScriptRoot

Write-Host "Verificando conexión a base de datos..." -ForegroundColor Yellow

# Ejecutar el servidor
Write-Host "Servidor iniciado en http://localhost:5000" -ForegroundColor Green
Write-Host "Presiona Ctrl+C para detener el servidor" -ForegroundColor Yellow
Write-Host ""

go run cmd/api/main.go
