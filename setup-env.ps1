# Script para configurar variables de entorno
# Ejecutar: .\setup-env.ps1

Write-Host "Configurando variables de entorno para desarrollo local..." -ForegroundColor Green

# Verificar si existe env.example
if (-not (Test-Path "env.example")) {
    Write-Host "Error: No se encontró el archivo env.example" -ForegroundColor Red
    exit 1
}

# Copiar env.example a .env
if (Test-Path ".env") {
    Write-Host "El archivo .env ya existe. ¿Deseas sobrescribirlo? (y/N)" -ForegroundColor Yellow
    $response = Read-Host
    if ($response -ne "y" -and $response -ne "Y") {
        Write-Host "Operación cancelada." -ForegroundColor Yellow
        exit 0
    }
}

Copy-Item "env.example" ".env"
Write-Host "Archivo .env creado exitosamente desde env.example" -ForegroundColor Green

Write-Host "`nVariables de entorno configuradas:" -ForegroundColor Cyan
Write-Host "- PORT=8080" -ForegroundColor White
Write-Host "- GIN_MODE=debug" -ForegroundColor White
Write-Host "- JWT_SECRET=crabi-jwt-secret-key-for-development-only" -ForegroundColor White
Write-Host "- DB_PATH=./data/crabi.db" -ForegroundColor White
Write-Host "- PLD_SERVICE_URL=http://98.81.235.22" -ForegroundColor White

Write-Host "`nPara desarrollo local, puedes editar el archivo .env según tus necesidades." -ForegroundColor Yellow
Write-Host "¡Listo para ejecutar el proyecto!" -ForegroundColor Green 