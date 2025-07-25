# 🐳 Docker - Crabi API

Guía completa para ejecutar el proyecto con Docker.

## 🚀 Quick Start

### Requisitos Previos

- **Docker Desktop** instalado y ejecutándose
- **Docker Compose** (incluido en Docker Desktop)

### Ejecutar Proyecto

```bash
# 1. Configurar variables de entorno
.\setup-env.ps1

# 2. Construir y ejecutar
docker-compose up --build
```

## 📁 Estructura Docker

```
crabi-test/
├── Dockerfile                    # Multi-stage build
├── docker-compose.yml            # Orquestación
├── .dockerignore                 # Archivos a ignorar
└── docker/
    └── pld-service/              # Servicio PLD simulado
        ├── index.html            # Página de documentación
        └── nginx.conf            # Configuración nginx
```

## 🔧 Dockerfile

### Multi-Stage Build

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

# Instalar dependencias
RUN apk add --no-cache git gcc musl-dev

# Configurar directorio de trabajo
WORKDIR /app

# Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Generar documentación Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/server/main.go

# Construir aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Instalar runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

# Crear usuario no-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Configurar directorio
WORKDIR /root/

# Copiar binario
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# Cambiar propietario
RUN chown -R appuser:appgroup /root/

# Usar usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Comando de inicio
CMD ["./main"]
```

## 🐙 Docker Compose

### Servicios

```yaml
services:
  crabi-api:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
    environment:
      - DOCKER_ENV=true
      - DB_PATH=/data/crabi.db
      - PLD_SERVICE_URL=http://pld.com
    volumes:
      - ./data:/data
    depends_on:
      - pld-service
    networks:
      - crabi-network
    restart: unless-stopped

  pld-service:
    image: nginx:alpine
    ports:
      - "3000:80"
    volumes:
      - ./docker/pld-service:/usr/share/nginx/html
      - ./docker/pld-service/nginx.conf:/etc/nginx/nginx.conf
    networks:
      - crabi-network
    restart: unless-stopped
```

## 📊 Comandos Docker

### Básicos

```bash
# Construir imagen
docker-compose build

# Ejecutar servicios
docker-compose up

# Ejecutar en background
docker-compose up -d

# Ver logs
docker-compose logs -f

# Ver logs de servicio específico
docker-compose logs -f crabi-api

# Detener servicios
docker-compose down

# Detener y remover volúmenes
docker-compose down -v
```

### Desarrollo

```bash
# Reconstruir sin cache
docker-compose build --no-cache

# Ejecutar con rebuild automático
docker-compose up --build

# Ejecutar solo un servicio
docker-compose up crabi-api

# Entrar al contenedor
docker-compose exec crabi-api sh

# Ver estado de servicios
docker-compose ps
```

### Debugging

```bash
# Ver logs detallados
docker-compose logs -f --tail=100

# Ver recursos utilizados
docker stats

# Ver información de imagen
docker inspect crabi-test-crabi-api

# Ver variables de entorno
docker-compose exec crabi-api env
```

## 🔍 Puertos y URLs

| Servicio | Puerto | URL | Descripción |
|----------|--------|-----|-------------|
| **crabi-api** | 8080 | http://localhost:8080 | API principal |
| **Swagger UI** | 8080 | http://localhost:8080/swagger/index.html | Documentación |
| **Health Check** | 8080 | http://localhost:8080/health | Estado de API |
| **pld-service** | 3000 | http://localhost:3000 | Servicio PLD simulado |

## 🔧 Variables de Entorno

### Docker Environment

```ini
# Configuración Docker
DOCKER_ENV=true

# Servidor
PORT=8080
GIN_MODE=debug

# JWT
JWT_SECRET=ITWVUfKiPTVWiVFvgjCYaOip6EejNAStO9+R5EbMM84=

# Base de datos
DB_PATH=/data/crabi.db

# Servicio PLD
PLD_SERVICE_URL=http://98.81.235.22
```

### Cargar desde .env

```yaml
env_file:
  - .env
```

## 🚨 Troubleshooting

### Error: "Docker Desktop not running"
```bash
# Iniciar Docker Desktop
# En Windows: Buscar "Docker Desktop" y ejecutar
```

### Error: "Port already in use"
```bash
# Verificar puertos en uso
netstat -ano | findstr :8080

# Cambiar puerto en docker-compose.yml
ports:
  - "8081:8080"  # Cambiar 8080 por 8081
```

### Error: "Build failed"
```bash
# Limpiar cache
docker system prune -a

# Reconstruir sin cache
docker-compose build --no-cache
```

### Error: "Permission denied"
```bash
# En Linux/Mac, dar permisos
chmod +x setup-env.ps1
chmod +x generate-jwt-secret.ps1
```

### Error: "SQLite CGO"
```bash
# Ya resuelto usando modernc.org/sqlite
# No requiere CGO
```

## 📈 Métricas Docker

### Tamaño de Imagen
- **Imagen final**: ~15MB
- **Imagen builder**: ~500MB
- **Optimización**: Multi-stage build

### Tiempo de Build
- **Primera vez**: ~60s
- **Con cache**: ~20s
- **Sin cache**: ~60s

### Recursos
- **CPU**: ~0.5 cores
- **RAM**: ~50MB
- **Disco**: ~100MB

## 🔄 CI/CD

### GitHub Actions (ejemplo)

```yaml
name: Docker Build

on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Build Docker image
      run: docker-compose build
    
    - name: Run tests
      run: docker-compose run crabi-api go test ./...
    
    - name: Push to registry
      run: docker push your-registry/crabi-api
```

## 🛠️ Optimizaciones

### Multi-Stage Build
- Reduce tamaño de imagen final
- Separa build de runtime
- Mejora seguridad

### Usuario No-Root
- Mejora seguridad
- Evita problemas de permisos
- Buenas prácticas

### Cache de Dependencias
- Acelera builds
- Reduce ancho de banda
- Mejora desarrollo

## 📝 Logs

### Ver Logs en Tiempo Real
```bash
docker-compose logs -f crabi-api
```

### Logs Típicos
```
crabi-api-1    | 2025/07/25 08:51:34 Base de datos SQLite inicializada correctamente
crabi-api-1    | 2025/07/25 08:51:34 Servidor iniciando en puerto 8080
crabi-api-1    | 2025/07/25 08:51:34 Documentación disponible en: http://localhost:8080/swagger/index.html
```

## 🎯 Próximos Pasos

- [ ] Configurar Docker Registry
- [ ] Implementar health checks
- [ ] Configurar monitoring
- [ ] Optimizar imagen base
- [ ] Implementar secrets management

---

**¡Docker listo y funcionando! 🚀** 