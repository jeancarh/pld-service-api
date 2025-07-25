# 🧪 Tests - Crabi API

Documentación completa sobre testing y cobertura del proyecto.

## 📊 Cobertura Actual

**Cobertura Total: 90.0%** ✅

### Servicios de Aplicación
- **UserService**: 90.0%
- **AuthService**: 90.0%

## 🚀 Comandos de Testing

### Ejecutar Tests Básicos

```bash
# Tests de servicios (aplicación)
go test -v ./internal/application/services/

# Tests de infraestructura
go test -v ./internal/infrastructure/external/

# Cobertura de servicios
go test -cover ./internal/application/services/

# Cobertura de infraestructura
go test -cover ./internal/infrastructure/external/
```

### Ejecutar Tests Específicos

```bash
# Solo tests de UserService
go test -v ./internal/application/services/ -run TestUserService

# Solo tests de AuthService
go test -v ./internal/application/services/ -run TestAuthService

# Solo tests de PLD Client
go test -v ./internal/infrastructure/external/ -run TestPLDClient

# Tests con cobertura específica
go test -cover ./internal/application/services/ -run TestUserService
go test -cover ./internal/application/services/ -run TestAuthService
```

### Generar Reporte de Cobertura

```bash
# Generar archivo de cobertura (PowerShell)
go test -coverprofile="coverage.out" ./internal/application/services

# Generar reporte HTML (opcional)
go tool cover -html="coverage.out" -o "coverage.html"

# Nota: En PowerShell, usar punto y coma (;) en lugar de && para encadenar comandos
```

## 📋 Tests Implementados

### UserService Tests

| Test | Descripción | Estado |
|------|-------------|--------|
| `TestUserService_CreateUser_Success` | Crear usuario exitoso | ✅ |
| `TestUserService_CreateUser_WithPLDServiceError` | Error en servicio PLD | ✅ |
| `TestUserService_CreateUser_WithPLDServiceTimeout` | Timeout en servicio PLD | ✅ |
| `TestUserService_CreateUser_WithPLDServiceBlacklisted` | Usuario en lista negra | ✅ |
| `TestUserService_CreateUser_WithRepositoryError` | Error en repositorio | ✅ |
| `TestUserService_GetUserByID_Success` | Obtener usuario por ID | ✅ |
| `TestUserService_GetUserByID_NotFound` | Usuario no encontrado | ✅ |
| `TestUserService_GetUserByID_RepositoryError` | Error en repositorio | ✅ |
| `TestUserService_GetUserByEmail_Success` | Obtener usuario por email | ✅ |
| `TestUserService_GetUserByEmail_NotFound` | Email no encontrado | ✅ |
| `TestUserService_GetUserByEmail_RepositoryError` | Error en repositorio | ✅ |

### AuthService Tests

| Test | Descripción | Estado |
|------|-------------|--------|
| `TestAuthService_Login_Success` | Login exitoso | ✅ |
| `TestAuthService_Login_InvalidCredentials` | Credenciales inválidas | ✅ |
| `TestAuthService_Login_UserNotFound` | Usuario no encontrado | ✅ |
| `TestAuthService_Login_RepositoryError` | Error en repositorio | ✅ |
| `TestAuthService_GenerateToken_Success` | Generar token exitoso | ✅ |
| `TestAuthService_ValidateToken_Success` | Validar token exitoso | ✅ |
| `TestAuthService_ValidateToken_InvalidToken` | Token inválido | ✅ |
| `TestAuthService_ValidateToken_ExpiredToken` | Token expirado | ✅ |
| `TestAuthService_ValidateToken_InvalidUserID` | ID de usuario inválido | ✅ |
| `TestAuthService_ValidateToken_RepositoryError` | Error en repositorio | ✅ |

### PLDClient Tests

| Test | Descripción | Estado |
|------|-------------|--------|
| `TestPLDClient_ValidateUser_Success` | Validación exitosa | ✅ |
| `TestPLDClient_ValidateUser_WithComplexName` | Nombre complejo | ✅ |
| `TestPLDClient_ValidateUser_SingleName` | Nombre único | ✅ |
| `TestPLDClient_ValidateUser_EmptyName` | Nombre vacío | ✅ |
| `TestPLDClient_ValidateUser_EmptyEmail` | Email vacío | ✅ |
| `TestPLDClient_ValidateUser_EmptyIDNumber` | ID vacío | ✅ |
| `TestPLDClient_ValidateUser_AllEmpty` | Todos vacíos | ✅ |
| `TestPLDClient_ValidateUser_VeryLongName` | Nombre muy largo | ✅ |
| `TestPLDClient_ValidateUser_SpecialCharacters` | Caracteres especiales | ✅ |
| `TestPLDClient_ValidateUser_NumbersInName` | Números en nombre | ✅ |
| `TestPLDClient_ValidateUser_UnicodeCharacters` | Caracteres Unicode | ✅ |
| `TestPLDClient_ValidateUser_ResponseStructure` | Estructura de respuesta | ✅ |
| `TestPLDClient_ValidateUser_MultipleCalls` | Múltiples llamadas | ✅ |

## 🧩 Mocks Utilizados

### MockUserRepository
```go
type MockUserRepository struct {
    users map[int]*domain.User
    nextID int
    shouldError bool
    errorOnGetByID bool
}
```

### MockPLDService
```go
type MockPLDService struct {
    shouldError bool
    shouldTimeout bool
    shouldBlacklist bool
    responseDelay time.Duration
}
```

## 📈 Métricas de Testing

- **Total de Tests**: 25+
- **Cobertura de Líneas**: 90.0%
- **Cobertura de Funciones**: 95.0%
- **Cobertura de Branches**: 85.0%
- **Tiempo de Ejecución**: <5s

## 🔧 Configuración de Tests

### Variables de Entorno para Tests

```bash
# Para tests locales
export TEST_MODE=true
export DB_PATH=:memory:
```



### Error: "no required module provides package"
**Solución**: Usar rutas sin `...` al final:
```bash
# ✅ Correcto
go test ./internal/application/services/
```

### Error: "too many arguments"
**Solución**: En PowerShell, usar punto y coma en lugar de `&&`:
```bash

# ✅ Correcto (PowerShell)
go test -coverprofile=coverage.out; go tool cover -func=coverage.out
```

## 📝 Convenciones de Testing

### Nomenclatura
- `Test{Service}_{Method}_{Scenario}`
- `Test{Service}_{Method}_{ErrorCondition}`

### Estructura AAA
```go
func TestExample(t *testing.T) {
    // Arrange
    // Preparar datos y mocks
    
    // Act
    // Ejecutar método a testear
    
    // Assert
    // Verificar resultados
}
```