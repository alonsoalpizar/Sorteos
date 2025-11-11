# Frontend builder stage
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

# Copiar package files del frontend
COPY frontend/package*.json ./

# Instalar dependencias
RUN npm ci

# Copiar código fuente del frontend
COPY frontend/ ./

# Build frontend
RUN npm run build

# Backend builder stage
FROM golang:1.22-alpine AS backend-builder

# Instalar dependencias de compilación
RUN apk add --no-cache git make

# Directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum
COPY backend/go.mod backend/go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY backend/ .

# Actualizar dependencias y generar go.sum
RUN go mod tidy

# Compilar binario
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Runtime stage
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copiar binario desde builder
COPY --from=backend-builder /app/main .

# Copiar archivos necesarios
COPY --from=backend-builder /app/migrations ./migrations

# Copiar frontend build
COPY --from=frontend-builder /frontend/dist ./frontend/dist

# Cambiar ownership
RUN chown -R appuser:appuser /app

# Cambiar a usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando por defecto
CMD ["./main"]
