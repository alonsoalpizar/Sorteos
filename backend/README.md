# Sorteos Platform - Backend API

Backend de la Plataforma de Sorteos construido con Go, Gin, PostgreSQL y Redis.

## 📋 Requisitos

- Go 1.22 o superior
- Docker y Docker Compose
- PostgreSQL 15+
- Redis 7+
- golang-migrate CLI (para migraciones)

## 🚀 Quick Start

### 1. Clonar repositorio y configurar entorno

```bash
cd /opt/Sorteos/backend
cp .env.example .env
# Editar .env con tus configuraciones
```

### 2. Iniciar servicios con Docker

```bash
# Desde la raíz del proyecto (/opt/Sorteos)
docker compose up -d postgres redis

# Verificar que los servicios están corriendo
docker compose ps

# Ver logs
docker compose logs -f postgres redis
```

### 3. Ejecutar migraciones

```bash
# Aplicar todas las migraciones
make migrate-up

# O manualmente:
migrate -path ./migrations \
  -database "postgresql://sorteos_user:sorteos_password@localhost:5432/sorteos_db?sslmode=disable" \
  up
```

### 4. Iniciar servidor de desarrollo

```bash
# Opción 1: Con make
make run

# Opción 2: Directamente con Go
go run cmd/api/main.go

# Opción 3: Con hot reload (requiere air)
go install github.com/cosmtrek/air@latest
air
```

El servidor estará disponible en: `http://localhost:8080`

## 🧪 Testing

### Health Checks

```bash
# Health check simple
curl http://localhost:8080/health

# Readiness check (verifica dependencias)
curl http://localhost:8080/ready

# Ping endpoint
curl http://localhost:8080/api/v1/ping
```

### Verificar PostgreSQL

```bash
# Conectar con psql
docker compose exec postgres psql -U sorteos_user -d sorteos_db

# Listar tablas
\dt

# Ver estructura de tabla users
\d users

# Contar usuarios
SELECT COUNT(*) FROM users;

# Salir
\q
```

### Verificar Redis

```bash
# Conectar con redis-cli
docker compose exec redis redis-cli

# Probar conexión
PING

# Ver todas las keys
KEYS *

# Salir
exit
```

## 📁 Estructura del Proyecto

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/                  # Entidades de dominio
│   │   ├── user.go
│   │   ├── raffle.go
│   │   ├── reservation.go
│   │   └── payment.go
│   ├── usecase/                 # Lógica de negocio
│   │   ├── auth/
│   │   ├── raffle/
│   │   ├── reservation/
│   │   └── payment/
│   └── adapters/                # Adaptadores externos
│       ├── http/                # Handlers HTTP
│       ├── db/                  # Repositorios PostgreSQL
│       ├── redis/               # Cliente Redis
│       ├── notifier/            # Emails/SMS
│       └── payments/            # Stripe
├── pkg/                         # Utilidades compartidas
│   ├── config/                  # Configuración
│   ├── logger/                  # Logger (Zap)
│   └── errors/                  # Errores personalizados
├── migrations/                  # Migraciones SQL
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_user_consents_table.up.sql
│   ├── 002_create_user_consents_table.down.sql
│   ├── 003_create_audit_logs_table.up.sql
│   └── 003_create_audit_logs_table.down.sql
├── .env.example                 # Variables de entorno
├── .gitignore
├── Dockerfile
├── Makefile                     # Comandos útiles
├── go.mod
└── README.md
```

## 🛠️ Comandos Útiles (Makefile)

```bash
make help             # Mostrar ayuda
make run              # Ejecutar aplicación
make build            # Compilar binario
make test             # Ejecutar tests
make test-coverage    # Tests con coverage
make migrate-up       # Aplicar migraciones
make migrate-down     # Revertir última migración
make migrate-create NAME=nombre  # Crear nueva migración
make docker-up        # Levantar contenedores
make docker-down      # Detener contenedores
make docker-logs      # Ver logs de contenedores
make lint             # Ejecutar linter
make clean            # Limpiar archivos build
```

## 🔐 Variables de Entorno Críticas

### Base de Datos (PostgreSQL)
```bash
CONFIG_DB_HOST=localhost
CONFIG_DB_PORT=5432
CONFIG_DB_USER=sorteos_user
CONFIG_DB_PASSWORD=sorteos_password
CONFIG_DB_NAME=sorteos_db
```

### Redis
```bash
CONFIG_REDIS_HOST=localhost
CONFIG_REDIS_PORT=6379
CONFIG_REDIS_PASSWORD=
```

### JWT (¡CAMBIAR EN PRODUCCIÓN!)
```bash
CONFIG_JWT_SECRET=change-this-to-a-secure-random-string-min-32-chars
CONFIG_JWT_ACCESS_TOKEN_EXPIRY=15m
CONFIG_JWT_REFRESH_TOKEN_EXPIRY=168h
```

### Stripe (Obtener de https://dashboard.stripe.com)
```bash
CONFIG_STRIPE_SECRET_KEY=sk_test_your_key_here
CONFIG_STRIPE_WEBHOOK_SECRET=whsec_your_secret_here
```

### SendGrid (Para emails)
```bash
CONFIG_SENDGRID_API_KEY=your_api_key_here
CONFIG_SENDGRID_FROM_EMAIL=noreply@sorteos.com
```

## 📊 Migraciones

### Crear nueva migración

```bash
make migrate-create NAME=add_raffles_table

# O manualmente:
migrate create -ext sql -dir migrations -seq add_raffles_table
```

Esto crea dos archivos:
- `XXX_add_raffles_table.up.sql` - Aplicar cambio
- `XXX_add_raffles_table.down.sql` - Revertir cambio

### Aplicar migraciones

```bash
# Aplicar todas pendientes
make migrate-up

# Aplicar N migraciones
migrate -path ./migrations -database "$DB_URL" up 2

# Ver versión actual
migrate -path ./migrations -database "$DB_URL" version
```

### Revertir migraciones

```bash
# Revertir última
make migrate-down

# Revertir N migraciones
migrate -path ./migrations -database "$DB_URL" down 2

# Revertir todas
migrate -path ./migrations -database "$DB_URL" down -all
```

### Forzar versión (¡CUIDADO!)

```bash
# Si una migración falla y la DB queda en estado inconsistente
migrate -path ./migrations -database "$DB_URL" force VERSION
```

## 🐛 Debugging

### Logs con diferentes niveles

```bash
# Development (logs detallados)
CONFIG_ENVIRONMENT=development go run cmd/api/main.go

# Debug (incluye queries SQL)
LOG_LEVEL=debug go run cmd/api/main.go

# Production (solo errores)
CONFIG_ENVIRONMENT=production go run cmd/api/main.go
```

### Adminer (UI para PostgreSQL)

```bash
# Iniciar con perfil debug
docker compose --profile debug up -d adminer

# Acceder a: http://localhost:8082
# Server: postgres
# User: sorteos_user
# Password: sorteos_password
# Database: sorteos_db
```

### Redis Commander (UI para Redis)

```bash
# Iniciar con perfil debug
docker compose --profile debug up -d redis-commander

# Acceder a: http://localhost:8081
```

## 🔒 Seguridad

### Configuración de Producción

1. **JWT Secret**: Mínimo 32 caracteres aleatorios
   ```bash
   openssl rand -base64 32
   ```

2. **PostgreSQL**:
   - Cambiar password
   - Habilitar SSL: `CONFIG_DB_SSLMODE=require`

3. **Redis**:
   - Configurar password
   - Habilitar TLS

4. **CORS**:
   - Configurar `CONFIG_ALLOWED_ORIGINS` con dominios específicos
   - Nunca usar `*` en producción

5. **Rate Limiting**:
   - Ajustar según tráfico esperado
   - Monitorear logs de rate limit

## 🚢 Deployment con Docker

### Build de imagen

```bash
cd /opt/Sorteos
docker compose build api
```

### Ejecutar en producción

```bash
# Configurar .env para producción
CONFIG_ENVIRONMENT=production

# Levantar todos los servicios
docker compose up -d

# Ver logs
docker compose logs -f api

# Escalar API (múltiples instancias)
docker compose up -d --scale api=3
```

### Health Checks

Docker verificará automáticamente la salud del contenedor:
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

## 📈 Monitoreo (Fase futura)

```bash
# Iniciar Prometheus + Grafana
docker compose --profile monitoring up -d prometheus grafana

# Prometheus: http://localhost:9090
# Grafana: http://localhost:3000 (admin/admin)
```

## 🧪 Tests

```bash
# Ejecutar todos los tests
make test

# Con coverage
make test-coverage

# Tests de integración (requiere Docker)
make test-integration

# Test específico
go test -v ./internal/usecase/auth/...

# Test con race detector
go test -race ./...
```

## 📝 Notas de Desarrollo

### Hot Reload con Air

Instalar Air:
```bash
go install github.com/cosmtrek/air@latest
```

Crear `.air.toml`:
```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/api"
bin = "tmp/main"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]
```

Ejecutar:
```bash
air
```

### Generación de código (opcional)

```bash
# Instalar herramientas
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación Swagger
swag init -g cmd/api/main.go

# Linter
golangci-lint run
```

## 🆘 Troubleshooting

### Error: "database connection failed"

```bash
# Verificar que PostgreSQL está corriendo
docker compose ps postgres

# Ver logs
docker compose logs postgres

# Reiniciar servicio
docker compose restart postgres
```

### Error: "redis connection failed"

```bash
# Verificar Redis
docker compose exec redis redis-cli PING

# Ver logs
docker compose logs redis

# Limpiar datos (¡CUIDADO!)
docker compose exec redis redis-cli FLUSHALL
```

### Error: "port 8080 already in use"

```bash
# Ver qué proceso usa el puerto
sudo lsof -i :8080

# Matar proceso
kill -9 PID

# O cambiar puerto en .env
CONFIG_SERVER_PORT=8081
```

### Limpiar todo y empezar de nuevo

```bash
# Detener contenedores
docker compose down

# Eliminar volúmenes (¡BORRA DATOS!)
docker compose down -v

# Limpiar imágenes
docker compose down --rmi all

# Reconstruir
docker compose up -d --build
```

## 📚 Referencias

- [Documentación completa](../Documentacion/README.md)
- [CLAUDE.md](../CLAUDE.md) - Contexto rápido para AI
- [Roadmap](../Documentacion/roadmap.md) - Plan de desarrollo
- [Stack Técnico](../Documentacion/stack_tecnico.md)
- [Seguridad](../Documentacion/seguridad.md)

## 📄 Licencia

Propietario: Ing. Alonso Alpízar
Fecha: Noviembre 2025
