# Build stage
FROM golang:1.23-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git gcc musl-dev

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Generar documentación Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/server/main.go

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Instalar dependencias de runtime
RUN apk --no-cache add ca-certificates sqlite

# Crear usuario no-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Establecer directorio de trabajo
WORKDIR /root/

# Copiar binario desde el stage de build
COPY --from=builder /app/main .

# Copiar documentación Swagger
COPY --from=builder /app/docs ./docs

# Cambiar propietario de archivos
RUN chown -R appuser:appgroup /root/

# Cambiar a usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"] 