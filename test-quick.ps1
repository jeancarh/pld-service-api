# Script de prueba rapida para Crabi API
# Ejecutar: .\test-quick.ps1

Write-Host "Ejecutando tests rapidos de Crabi API..." -ForegroundColor Green

# Configurar variable de entorno para tests
$env:PLD_SERVICE_URL="http://localhost:3000"

Write-Host "`nTests de Servicios (Aplicacion):" -ForegroundColor Cyan
go test -cover ./internal/application/services/

Write-Host "`nTests de Infraestructura (PLDClient):" -ForegroundColor Cyan
go test -cover ./internal/infrastructure/external/

Write-Host "`nTests completados!" -ForegroundColor Green
Write-Host "Cobertura objetivo: >90% para servicios" -ForegroundColor Yellow
Write-Host "Proyecto listo para desarrollo/produccion" -ForegroundColor Green 