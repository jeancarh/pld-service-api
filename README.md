# Crabi API - Sistema de PLD

Sistema de gestión de usuarios con validación PLD (Prevención de Lavado de Dinero) implementado en Go con arquitectura hexagonal.

## 📋 Características

- ✅ **Arquitectura Hexagonal** (Clean Architecture)
- ✅ **Autenticación JWT**
- ✅ **Validación PLD** con servicio externo
- ✅ **Base de datos SQLite** (local) / **modernc.org/sqlite** (sin CGO)
- ✅ **Documentación OpenAPI/Swagger** automática
- ✅ **Tests unitarios** con cobertura >90%
- ✅ **Docker & Docker Compose**
- ✅ **Validadores personalizados**
- ✅ **Principios SOLID** aplicados

## 📋 Prerrequisitos

### Para Docker (Recomendado)
- ✅ **Docker Desktop** instalado y ejecutándose
- ✅ **Git** para clonar el repositorio

### Para Desarrollo Local
- ✅ **Go 1.23+** instalado
- ✅ **Git** para clonar el repositorio
- ✅ **PowerShell** (para scripts de configuración)

### Verificar instalaciones:
```bash
# Verificar Go
go version

# Verificar Docker
docker --version
docker-compose --version

# Verificar Git
git --version
```

## 🚀 Quick Start

### Opción 1: Docker (Recomendado)

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd crabi-test

# 2. Configurar variables de entorno
.\setup-env.ps1

# 3. Ejecutar con Docker Compose
docker-compose up --build
```

### Opción 2: Desarrollo Local

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd crabi-test

# 2. Configurar variables de entorno
.\setup-env.ps1

# 3. Gestionar dependencias
go mod tidy
go mod download

# 4. Generar documentación Swagger
swag init -g cmd/server/main.go

# 5. Ejecutar aplicación
go run cmd/server/main.go
```

## 🌐 URLs de Acceso

Una vez ejecutada la aplicación:

- **API Base**: http://localhost:8080
- **Documentación Swagger**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health
- **Colección Postman**: Importar `Crabi_API.postman_collection.json`

## 📦 Gestión de Dependencias

### Comandos básicos

```bash
# Descargar dependencias
go mod download

# Limpiar y organizar dependencias
go mod tidy

# Ver dependencias
go mod graph

# Verificar dependencias
go mod verify
```

### Agregar nuevas dependencias

```bash
# Agregar dependencia
go get github.com/nueva/dependencia

# Agregar dependencia específica
go get github.com/nueva/dependencia@v1.2.3

# Actualizar dependencia
go get -u github.com/existente/dependencia
```

### Dependencias actuales

El proyecto usa las siguientes dependencias principales:

- **Gin**: Framework web HTTP
- **JWT**: Autenticación con tokens
- **Validator**: Validación de datos
- **SQLite**: Base de datos local
- **Swagger**: Documentación automática
- **Godotenv**: Variables de entorno

### Archivos de dependencias

- **go.mod**: Define módulo y dependencias directas
- **go.sum**: Checksums de seguridad de dependencias

## 🔧 Troubleshooting

### Problemas Comunes

#### Docker no inicia
```bash
# Verificar que Docker Desktop esté ejecutándose
# Verificar puertos disponibles
netstat -an | findstr :8080
```

#### Error de variables de entorno
```bash
# Regenerar archivo .env
.\setup-env.ps1

# Verificar que el archivo existe
ls .env
```

#### Error de dependencias Go
```bash
# Limpiar cache de módulos
go clean -modcache

# Descargar dependencias nuevamente
go mod download

# Verificar y limpiar dependencias
go mod tidy
go mod verify

# Forzar descarga de dependencias
go mod download -x
```

#### Error de permisos PowerShell
```bash
# Ejecutar PowerShell como administrador
# O cambiar política de ejecución
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Logs y Debugging

```bash
# Ver logs de Docker
docker-compose logs -f

# Ver logs de la aplicación
docker-compose logs crabi-api

# Ejecutar en modo debug
GIN_MODE=debug go run cmd/server/main.go
```

## 📊 Tests y Cobertura

### Ejecutar Tests

```bash
# Tests de servicios (aplicación)
go test -v ./internal/application/services/

# Tests de infraestructura
go test -v ./internal/infrastructure/external/

# Cobertura de servicios
go test -cover ./internal/application/services/

# Cobertura de infraestructura
go test -cover ./internal/infrastructure/external/

# Script de prueba rápida (recomendado)
.\test-quick.ps1
```

### Cobertura Actual: **90.0%** ✅

```
# Servicios (Aplicación)
PASS
coverage: 90.0% of statements
ok      crabi-test/internal/application/services

# Infraestructura (PLDClient)
PASS
coverage: 75.0% of statements
ok      crabi-test/internal/infrastructure/external
```

## 🔧 Configuración

### Variables de Entorno

Crear archivo `.env` basado en `env.example`:

```ini
# Configuración del servidor
PORT=8080
GIN_MODE=debug

JWT_SECRET=ITWVUfKiPTVWiVFvgjCYaOip6EejNAStO9+R5EbMM84=

# Base de datos
DB_PATH=./data/crabi.db

# Servicio PLD (URL real)
PLD_SERVICE_URL=http://98.81.235.22

# Docker environment
DOCKER_ENV=true
```


## 📚 Documentación Swagger

### Generar Documentación

La documentación se genera automáticamente desde los comentarios en el código:

```bash
# Instalar swag (si no está instalado)
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
swag init -g cmd/server/main.go

# Regenerar después de cambios
swag init -g cmd/server/main.go --parseDependency
```

### Anotaciones Swagger

Los endpoints están documentados con anotaciones:

```go
// CreateUser godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario validando contra el servicio PLD
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ...
}
```

### Acceder a la Documentación

- **URL**: http://localhost:8080/swagger/index.html
- **Disponible**: Cuando la aplicación está corriendo
- **Interactiva**: Puedes probar endpoints directamente

### Documentar DTOs

Los DTOs también se documentan con anotaciones:

```go
// CreateUserRequest representa la solicitud para crear un usuario
// @Description Solicitud para crear un nuevo usuario
type CreateUserRequest struct {
    // @Description Nombre completo del usuario
    // @Example "Juan Pérez"
    // @Required
    Name string `json:"name" binding:"required,min=2,max=100" example:"Juan Pérez"`
    
    // @Description Email del usuario (debe ser único)
    // @Example "juan.perez@email.com"
    // @Required
    Email string `json:"email" binding:"required,email" example:"juan.perez@email.com"`
}
```


## 📡 Endpoints

| Endpoint | Método | Descripción | Auth |
|----------|--------|-------------|------|
| `/health` | GET | Health check | ❌ |
| `/api/v1/users` | POST | Crear usuario | ❌ |
| `/api/v1/auth/login` | POST | Login | ❌ |
| `/api/v1/users/me` | GET | Usuario autenticado | ✅ |
| `/api/v1/users/:id` | GET | Usuario por ID | ✅ |
| `/api/v1/users/:id` | DELETE | Eliminar usuario | ✅ |
| `/swagger/index.html` | GET | Documentación | ❌ |

## 🧪 Testing

### Ejemplos de Requests

#### 1. Crear Usuario
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan.perez@email.com",
    "password": "password123",
    "id_number": "12345678"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan.perez@email.com",
    "password": "password123"
  }'
```

#### 3. Obtener Usuario (con token)
```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### 4. Eliminar Usuario (con token)
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🐳 Docker

### Estructura Docker

```
crabi-test/
├── Dockerfile              # Multi-stage build
├── docker-compose.yml      # Orchestration
└── docker/
    └── pld-service/        # Servicio PLD simulado
        ├── index.html
        └── nginx.conf
```

### Comandos Docker

```bash
# Construir y ejecutar
docker-compose up --build

# Solo construir
docker-compose build

# Ejecutar en background
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down
```

### Puertos

- **API**: `http://localhost:8080`
- **Swagger**: `http://localhost:8080/swagger/index.html`
- **PLD Service**: `http://98.81.235.22` (servicio real)

## 📁 Estructura del Proyecto

```
crabi-test/
├── cmd/
│   └── server/
│       └── main.go                 # Punto de entrada
├── internal/
│   ├── adapters/
│   │   └── repositories/           # Implementación de repositorios
│   ├── application/
│   │   ├── ports/                 # Interfaces (puertos)
│   │   └── services/              # Lógica de negocio
│   ├── domain/                    # Entidades de dominio
│   └── infrastructure/
│       ├── database/              # Configuración de BD
│       ├── external/              # Clientes externos (PLD)
│       └── http/                  # Handlers y middleware
├── pkg/
│   └── validator/                 # Validadores personalizados
├── tests/                         # Tests de integración
├── docker/                        # Configuración Docker
├── docs/                          # Documentación Swagger
├── data/                          # Base de datos SQLite
├── .env                           # Variables de entorno
├── env.example                    # Ejemplo de variables
├── setup-env.ps1                  # Script de configuración
├── docker-compose.yml             # Orquestación Docker
├── Dockerfile                     # Build de Docker
├── go.mod                         # Dependencias Go
├── go.sum                         # Checksums de dependencias
├── README.md                      # Este archivo
└── Crabi_API.postman_collection.json  # Colección Postman
```

## 🔍 Validaciones

### Crear Usuario
- ✅ Email válido
- ✅ Password mínimo 6 caracteres
- ✅ ID number numérico
- ✅ Usuario no en lista negra PLD
- ✅ Nombre no vacío

### Login
- ✅ Email válido
- ✅ Password correcto
- ✅ Usuario existe

### JWT Token
- ✅ Token válido
- ✅ Token no expirado
- ✅ Usuario existe en BD

## 🛠️ Tecnologías

- **Go 1.23** - Lenguaje principal
- **Gin** - Framework HTTP
- **modernc.org/sqlite** - Base de datos (sin CGO)
- **JWT** - Autenticación
- **Swagger** - Documentación API
- **Docker** - Containerización
- **bcrypt** - Hash de contraseñas
- **validator** - Validación de datos

## 🏗️ Principios SOLID y Patrones de Diseño

### ✅ **Principios SOLID Aplicados:**

#### **1. Single Responsibility Principle (SRP)**
- **UserService**: Responsable solo de la lógica de negocio de usuarios
- **AuthService**: Responsable solo de autenticación y JWT
- **PLDClient**: Responsable solo de comunicación con servicio PLD
- **UserRepository**: Responsable solo de persistencia de datos

#### **2. Open/Closed Principle (OCP)**
- **Interfaces**: `UserRepository`, `PLDService` permiten extensión sin modificación
- **Middleware**: Sistema de middleware extensible
- **Validators**: Validadores personalizados extensibles

#### **3. Liskov Substitution Principle (LSP)**
- **Mocks**: `MockUserRepository` sustituye `UserRepository` sin problemas
- **Mocks**: `MockPLDService` sustituye `PLDService` correctamente
- **Interfaces**: Todas las implementaciones son intercambiables

#### **4. Interface Segregation Principle (ISP)**
- **UserRepository**: Interface específica para operaciones de usuario
- **PLDService**: Interface específica para validación PLD
- **AuthService**: Interface específica para autenticación

#### **5. Dependency Inversion Principle (DIP)**
- **Dependency Injection**: Servicios dependen de interfaces, no implementaciones
- **Inversion of Control**: `main.go` configura las dependencias
- **Abstractions**: Interfaces definen contratos claros

### ✅ **Patrones de Diseño Implementados:**

#### **1. Repository Pattern**
```go
type UserRepository interface {
    Create(user *domain.User) error
    GetByID(id int) (*domain.User, error)
    GetByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
}
```
#### **2. Service Layer Pattern**
```go
type UserService struct {
    userRepo    repositories.UserRepository
    pldService  ports.PLDService
}
```

#### **3. Dependency Injection Pattern**
```go
func NewUserService(userRepo repositories.UserRepository, pldService ports.PLDService) *UserService {
    return &UserService{
        userRepo:   userRepo,
        pldService: pldService,
    }
}
```

#### **4. Factory Pattern**
```go
func NewPLDClient() *PLDClient {
    baseURL := os.Getenv("PLD_SERVICE_URL")
    return &PLDClient{baseURL: baseURL}
}
```

#### **5. Middleware Pattern**
```go
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
    }
}
```
- **Ventaja**: Funcionalidad transversal reutilizable
- **Beneficio**: Separación de concerns

#### **6. DTO Pattern (Data Transfer Object)**
```go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    IDNumber string `json:"id_number" binding:"required"`
}
```

#### **7. Hexagonal Architecture (Clean Architecture)**
```
┌─────────────────────────────────────┐
│           Infrastructure            │
│  ┌─────────────┐ ┌─────────────┐   │
│  │   HTTP      │ │  Database   │   │
│  │  Handlers   │ │  SQLite     │   │
│  └─────────────┘ └─────────────┘   │
└─────────────────────────────────────┘
           │           │
           ▼           ▼
┌─────────────────────────────────────┐
│           Application               │
│  ┌─────────────┐ ┌─────────────┐   │
│  │   Services  │ │   Ports     │   │
│  │  (Business) │ │(Interfaces) │   │
│  └─────────────┘ └─────────────┘   │
└─────────────────────────────────────┘
           │           │
           ▼           ▼
┌─────────────────────────────────────┐
│             Domain                  │
│  ┌─────────────┐ ┌─────────────┐   │
│  │   Entities  │ │   Value     │   │
│  │   (User)    │ │   Objects   │   │
│  └─────────────┘ └─────────────┘   │
└─────────────────────────────────────┘
```
- **Ventaja**: Separación clara de capas
- **Beneficio**: Independencia de frameworks

### ✅ **Beneficios de la Arquitectura:**

1. **Testabilidad**: Fácil mocking y testing unitario
2. **Mantenibilidad**: Código organizado y fácil de entender
3. **Escalabilidad**: Fácil agregar nuevas funcionalidades
4. **Flexibilidad**: Cambio de tecnologías sin afectar lógica de negocio
5. **Reutilización**: Componentes reutilizables en diferentes contextos

## 📈 Métricas

- **Cobertura de Tests**: 90.0%
- **Tiempo de Build**: ~60s (Docker)
- **Tamaño de Imagen**: ~15MB
- **Endpoints**: 6 totales
- **Validaciones**: 15+ reglas

## 🚨 Troubleshooting

### Error: "CGO_ENABLED=0, go-sqlite3 requires cgo"
**Solución**: Ya está resuelto usando `modernc.org/sqlite`

### Error: "PLD_SERVICE_URL no está configurado"
**Solución**: Ejecutar `.\setup-env.ps1` o configurar manualmente

### Error: "Docker Desktop not running"
**Solución**: Iniciar Docker Desktop

### Error: "Port 8080 already in use"
**Solución**: Cambiar `PORT` en `.env` o detener proceso existente


---

**¡Listo para usar! 🚀** 
