# Contexto del Proyecto para Claude AI

**Proyecto:** Plataforma de Sorteos/Rifas en Línea
**Propietario:** Ing. Alonso Alpízar
**Stack:** Go + React + PostgreSQL + Redis
**Estado:** Documentación técnica completa (100%)

---

## 🎯 Propósito de este Archivo

Este archivo proporciona contexto rápido a Claude AI para trabajar eficientemente en el proyecto sin necesidad de leer toda la documentación cada vez.

---

## 📋 Información Esencial

### Arquitectura

- **Backend:** Go 1.22+ con Gin (arquitectura hexagonal)
- **Frontend:** React 18 + TypeScript + Vite + Tailwind CSS + shadcn/ui
- **Base de Datos:** PostgreSQL 16 (transaccional) - **INSTALACIÓN LOCAL**
- **Cache/Locks:** Redis 7 (locks distribuidos, rate limiting) - **INSTALACIÓN LOCAL**
- **Pagos:** Stripe (MVP) → PayPal (Fase 2)
- **Servidor Web:** Nginx (reverse proxy + SSL)

### Estructura de Carpetas

```
/opt/Sorteos/
├── backend/
│   ├── cmd/api/              # Entry point
│   ├── internal/
│   │   ├── domain/           # Entidades (User, Raffle, Reservation, Payment)
│   │   ├── usecase/          # Casos de uso (CreateRaffle, ReserveNumbers, etc.)
│   │   └── adapters/         # HTTP, DB, Payments, Notifier
│   ├── pkg/                  # Logger, Config, Errors
│   ├── migrations/           # SQL migrations
│   ├── uploads/              # Archivos subidos
│   ├── .env                  # Variables de entorno
│   └── sorteos-api           # Binario compilado
├── frontend/
│   ├── src/
│   │   ├── app/              # Router, providers
│   │   ├── features/         # auth, raffles, checkout
│   │   ├── components/ui/    # shadcn/ui components
│   │   └── lib/              # Utilidades
│   ├── dist/                 # Build de producción (servido por backend)
│   └── public/
└── Documentacion/            # 10 documentos técnicos (181 KB)
```

---

## 🚨 RESTRICCIONES OBLIGATORIAS

### 1. Colores (CRÍTICO)

**⚠️ PROHIBIDO ABSOLUTAMENTE:**
- Morado, púrpura, violeta (#8B5CF6, #A855F7, etc.)
- Rosa, pink, magenta (#EC4899, #F472B6, etc.)
- Fucsia (#D946EF, #E879F9)
- Gradientes arcoíris o neón

**✅ PERMITIDO:**
- **Primary:** Azul #3B82F6 (confianza, profesionalismo)
- **Secondary:** Slate #64748B (elegancia, corporativo)
- **Neutral:** Grises #171717 a #FAFAFA
- **Success:** Verde #10B981 (solo confirmaciones)
- **Warning:** Ámbar #F59E0B (solo advertencias)
- **Error:** Rojo #EF4444 (solo errores)

**Referencias:** Stripe.com, Linear.app, Vercel.com, Coinbase.com

Ver: `Documentacion/estandar_visual.md` y `Documentacion/.paleta-visual-aprobada.md`

### 2. Seguridad

- **JWT:** Access token 15 min, Refresh token 7 días
- **Passwords:** bcrypt cost 12
- **Rate Limiting:** Redis (5-60 req/min según endpoint)
- **NUNCA:** Almacenar números de tarjeta (usar tokens de Stripe)
- **PCI DSS:** Delegado a Stripe
- **GDPR:** Derecho al olvido implementado

### 3. Concurrencia (CRÍTICO)

**Problema:** Doble venta de números de sorteo

**Solución obligatoria:**
1. Lock distribuido en Redis (SETNX con TTL)
2. Verificación en PostgreSQL (transacción)
3. Reserva temporal (5 min)

```go
// Ejemplo de lock distribuido
lockKey := fmt.Sprintf("lock:raffle:%d:num:%s", raffleID, number)
acquired := rdb.SetNX(ctx, lockKey, userID, 30*time.Second)
if !acquired {
    return errors.New("número ya reservado")
}
defer rdb.Del(ctx, lockKey)
```

Ver: `Documentacion/modulos.md` sección "Reserva y Compra"

---

## 📚 Documentación Disponible

1. **arquitecturaIdeaGeneral.md** - Visión general, concurrencia, flujos
2. **stack_tecnico.md** - Tecnologías, dependencias, versiones
3. **roadmap.md** - Fases de desarrollo (MVP → Fase 3)
4. **modulos.md** - 7 módulos con código y casos de uso
5. **estandar_visual.md** - Design system (Tailwind + shadcn/ui)
6. **seguridad.md** - JWT, RBAC, rate limiting, OWASP Top 10
7. **pagos_integraciones.md** - Stripe, webhooks, idempotencia
8. **parametrizacion_reglas.md** - Parámetros dinámicos (80+)
9. **operacion_backoffice.md** - Dashboard admin, liquidaciones
10. **terminos_y_condiciones_impacto.md** - GDPR, PCI DSS

---

## 🔑 Entidades Principales

### User
```go
type User struct {
    ID           int64
    Email        string
    Phone        string
    PasswordHash string
    Role         UserRole // user, admin
    KYCLevel     KYCLevel // none, email_verified, phone_verified, full_kyc
    Status       UserStatus
}
```

### Raffle
```go
type Raffle struct {
    ID            int64
    UserID        int64 // owner
    Title         string
    Status        RaffleStatus // draft, active, suspended, completed
    DrawDate      time.Time
    PricePerNumber decimal.Decimal
    TotalNumbers  int
    SoldCount     int
}
```

### Reservation
```go
type Reservation struct {
    ID             int64
    RaffleID       int64
    UserID         int64
    Numbers        []string
    Status         ReservationStatus // pending, confirmed, expired
    IdempotencyKey string
    ExpiresAt      time.Time
}
```

### Payment
```go
type Payment struct {
    ID             int64
    ReservationID  int64
    Provider       string // "stripe"
    Amount         decimal.Decimal
    Status         PaymentStatus
    ExternalID     string
    IdempotencyKey string
}
```

---

## 🛠️ Comandos Útiles - INSTALACIÓN LOCAL

### Backend (Go)
```bash
cd /opt/Sorteos/backend

# Compilar binario
go build -o sorteos-api ./cmd/api

# Reiniciar servicio
systemctl restart sorteos-api

# Ver logs
journalctl -xeu sorteos-api -f

# Ver estado
systemctl status sorteos-api
```

### Frontend (React)
```bash
cd /opt/Sorteos/frontend

# Desarrollo local
npm run dev           # Puerto 5173

# Build producción (10 segundos)
npm run build         # Output: dist/

# El backend ya sirve automáticamente desde dist/
```

### Servicios del Sistema
```bash
# PostgreSQL 16
systemctl status postgresql
psql -U sorteos_user -d sorteos_db

# Redis 7
systemctl status redis-server
redis-cli ping

# Backend API
systemctl status sorteos-api
systemctl restart sorteos-api
```

### Nginx (NO TOCAR - Ya configurado)
```bash
# Ver configuración
cat /etc/nginx/sites-available/sorteos.club

# Recargar (solo si es necesario)
nginx -t && systemctl reload nginx
```

---

## 🎨 Guía Rápida de UI

### Componentes Base (shadcn/ui)

```tsx
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'

// Button primary (azul)
<Button variant="default">Comprar Boleto</Button>

// Button secondary (slate/gris)
<Button variant="secondary">Ver Detalles</Button>

// Card profesional
<Card className="border-neutral-200">
  <CardHeader>
    <CardTitle className="text-neutral-900">Título</CardTitle>
  </CardHeader>
  <CardContent>
    Contenido
  </CardContent>
</Card>

// Badge de estado
<Badge className="bg-primary-100 text-primary-700">
  Activo
</Badge>
```

### Estados de Color

```tsx
// Success (verde)
<Alert className="bg-success/10 border-success/20 text-success">
  ✓ Operación exitosa
</Alert>

// Warning (ámbar)
<Alert className="bg-warning/10 border-warning/20 text-warning">
  ⚠ Advertencia
</Alert>

// Error (rojo)
<Alert className="bg-error/10 border-error/20 text-error">
  ✗ Error crítico
</Alert>
```

---

## 🔐 Endpoints Críticos

### Auth
- `POST /auth/register` - Registro con verificación email/SMS
- `POST /auth/login` - Login con JWT
- `POST /auth/refresh` - Refresh token

### Raffles
- `POST /raffles` - Crear sorteo (requiere KYC >= email_verified)
- `GET /raffles` - Listar (paginado, filtros)
- `GET /raffles/:id` - Detalle con números disponibles
- `POST /raffles/:id/publish` - Publicar (solo owner)

### Reservations (CRÍTICO - Alta concurrencia)
- `POST /raffles/:id/reservations` - **Reservar números con lock distribuido**
- `GET /reservations/:id` - Ver reserva
- `DELETE /reservations/:id` - Cancelar

### Payments
- `POST /payments` - Crear pago (idempotente con header `Idempotency-Key`)
- `POST /webhooks/stripe` - Webhook de Stripe (verificar firma)

### Admin
- `PATCH /admin/raffles/:id/suspend` - Suspender sorteo
- `PATCH /admin/users/:id/kyc` - Verificar KYC
- `POST /admin/settlements` - Crear liquidación

---

## ⚡ Flujos Críticos

### 1. Reserva y Compra de Números

```
Usuario → Selecciona números
       → POST /raffles/:id/reservations
          ├─ Lock Redis (SETNX) - 30s
          ├─ Verificar en DB (transacción)
          ├─ Crear reserva (expires_at = now + 5min)
          └─ Liberar lock
       → Frontend muestra timer 5 min
       → POST /payments (con Idempotency-Key)
          ├─ Stripe.js tokeniza tarjeta
          ├─ Backend crea PaymentIntent
          └─ Webhook confirma → marca números como sold
```

### 2. Publicación de Sorteo

```
Usuario → Crea sorteo (draft)
       → Sube imágenes
       → POST /raffles/:id/publish
          ├─ Validar parámetros (ver parametrizacion_reglas.md)
          ├─ Verificar KYC >= email_verified
          ├─ Validar max sorteos activos (default: 10)
          ├─ Generar números (00-99 o configurable)
          └─ status = active
```

### 3. Selección de Ganador

```
Cron job (diario a las 00:00)
  → Consultar Lotería Nacional CR API
  → Extraer número ganador (últimos 2 dígitos)
  → Buscar número en raffle_numbers
  → Si vendido:
      ├─ Actualizar raffle.winner_id
      ├─ Notificar ganador (email/SMS)
      └─ Crear settlement (calcular neto)
    Si no vendido:
      └─ raffle.winner_id = null
```

---

## 🧪 Tests Críticos

### Backend
```bash
# Test de concurrencia (reservas)
go test -v -race ./internal/usecase -run TestReserveNumbers_Concurrent

# Test de idempotencia (pagos)
go test -v ./internal/usecase -run TestCreatePayment_Idempotent
```

### Frontend
```bash
# Tests de componentes
npm run test

# Tests e2e (Playwright/Cypress)
npm run test:e2e
```

### Pruebas de Carga
```bash
# 1000 usuarios concurrentes comprando
k6 run scripts/load-test-reservations.js
```

**Criterio de éxito:** 0% de doble venta en 1000 requests concurrentes

---

## 📊 Métricas Clave

### Sistema
- `http_requests_total{method, path, status}` - Total requests
- `reservation_duration_seconds` - Latencia de reservas
- `payment_success_rate` - Tasa de éxito de pagos
- `active_reservations_gauge` - Reservas activas

### Negocio
- MAU (Monthly Active Users)
- Tasa de conversión reserva → pago (objetivo: 70%)
- Sorteos completados / mes
- Revenue total / comisiones

---

## ✅ Estado Actual del Sistema (2025-11-13)

### Infraestructura - MIGRACIÓN COMPLETADA: Docker → Local

**Antes (Docker):**
- 6 paquetes Docker + dependencias (464 MB overhead)
- Rebuild frontend: 3+ minutos
- Debugging complejo (logs en contenedores, exec, etc.)

**Ahora (Instalación Local):**
- ✅ PostgreSQL 16 nativo (puerto 5432)
- ✅ Redis 7 nativo (puerto 6379)
- ✅ Backend Go como servicio systemd (puerto 8080)
- ✅ Frontend servido por backend desde `frontend/dist/`
- ✅ Nginx como reverse proxy con SSL
- ✅ Rebuild frontend: **10 segundos** (`npm run build`)
- ✅ Logs centralizados en journalctl
- ✅ Debugging directo con herramientas estándar

### Servicios Activos

```bash
# Verificación de servicios
systemctl is-active postgresql redis-server sorteos-api nginx
# postgresql: active
# redis-server: active
# sorteos-api: active
# nginx: active
```

### Configuración de Servicios

**PostgreSQL:**
- Host: localhost:5432
- Database: sorteos_db
- User: sorteos_user
- Ubicación: /var/lib/postgresql/16/main

**Redis:**
- Host: localhost:6379
- Sin password
- Persistence: RDB + AOF

**Backend API (systemd):**
- Servicio: sorteos-api.service
- WorkingDirectory: /opt/Sorteos
- Binario: /opt/Sorteos/backend/sorteos-api
- Puerto: 8080
- Auto-start: enabled

**Nginx:**
- Proxy: https://sorteos.club → localhost:8080
- SSL: Configurado (certbot)
- Static files: Servidos por backend Go

---

## 🔧 Flujo de Trabajo de Desarrollo

### Cambios en Backend (Go)

```bash
cd /opt/Sorteos/backend

# 1. Hacer cambios en archivos .go

# 2. Compilar (verifica errores)
go build -o sorteos-api ./cmd/api

# 3. Reiniciar servicio
systemctl restart sorteos-api

# 4. Verificar logs
journalctl -xeu sorteos-api -f

# 5. Health check
curl http://localhost:8080/health
```

**Tiempo total:** ~5-10 segundos

### Cambios en Frontend (React/TypeScript)

```bash
cd /opt/Sorteos/frontend

# 1. Hacer cambios en archivos .tsx/.ts

# 2. Build (solo 10 segundos!)
npm run build

# 3. El backend ya está sirviendo el nuevo build
# No se requiere reiniciar nada

# 4. Verificar
curl -I https://sorteos.club/
```

**Tiempo total:** ~10 segundos (vs 3+ minutos con Docker)

### Cambios en Base de Datos

```bash
# Crear nueva migración
cd /opt/Sorteos/backend/migrations
touch 010_nueva_migracion.up.sql
touch 010_nueva_migracion.down.sql

# Aplicar migración (si usas migrate CLI)
migrate -path ./migrations -database "postgresql://sorteos_user:sorteos_password@localhost:5432/sorteos_db?sslmode=disable" up

# O aplicar manualmente
psql -U sorteos_user -d sorteos_db -f migrations/010_nueva_migracion.up.sql
```

### Verificaciones Post-Deploy

```bash
# 1. Servicios corriendo
systemctl is-active postgresql redis-server sorteos-api nginx

# 2. Health checks
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/ping

# 3. Frontend
curl -I https://sorteos.club/

# 4. API pública
curl https://sorteos.club/api/v1/ping

# 5. Logs sin errores
journalctl -xeu sorteos-api --since "5 minutes ago"
```

---

## 🔍 Debugging y Troubleshooting

### Ver Logs

```bash
# Backend API
journalctl -xeu sorteos-api -f
journalctl -xeu sorteos-api --since "1 hour ago"

# PostgreSQL
journalctl -xeu postgresql -f

# Redis
journalctl -xeu redis-server -f

# Nginx
tail -f /var/log/nginx/error.log
tail -f /var/log/nginx/access.log
```

### Conectar a Bases de Datos

```bash
# PostgreSQL
psql -U sorteos_user -d sorteos_db

# Redis
redis-cli
> PING
> KEYS *
> GET user:session:12345
```

### Verificar Puertos

```bash
# Ver qué está escuchando en cada puerto
ss -tlnp | grep -E "5432|6379|8080|80|443"

# PostgreSQL (5432)
# Redis (6379)
# Backend API (8080)
# Nginx (80, 443)
```

### Reiniciar Todo (Emergency)

```bash
# Reiniciar servicios en orden
systemctl restart postgresql
systemctl restart redis-server
systemctl restart sorteos-api
systemctl restart nginx

# Verificar estado
systemctl status postgresql redis-server sorteos-api nginx
```

---

## 📝 Variables de Entorno

Archivo: `/opt/Sorteos/backend/.env`

**Críticas:**
```bash
# Base de Datos
CONFIG_DB_HOST=localhost      # Era "postgres" en Docker
CONFIG_DB_PORT=5432
CONFIG_DB_USER=sorteos_user
CONFIG_DB_PASSWORD=sorteos_password
CONFIG_DB_NAME=sorteos_db

# Redis
CONFIG_REDIS_HOST=localhost   # Era "redis" en Docker
CONFIG_REDIS_PORT=6379
CONFIG_REDIS_PASSWORD=

# JWT
CONFIG_JWT_SECRET=change-this-to-a-secure-random-string-min-32-chars
CONFIG_JWT_ACCESS_TOKEN_EXPIRY=15m
CONFIG_JWT_REFRESH_TOKEN_EXPIRY=168h

# Server
CONFIG_ENV=development
CONFIG_PORT=8080

# Uploads
CONFIG_STORAGE_PATH=./backend/uploads  # Relativo a /opt/Sorteos
```

---

## 🆘 En Caso de Duda

1. **Colores:** Si no es azul/gris/verde/ámbar/rojo → NO USAR
2. **Concurrencia:** Siempre usar locks de Redis para reservas
3. **Pagos:** Siempre usar Idempotency-Key
4. **Seguridad:** Rate limiting en endpoints sensibles
5. **GDPR:** Nunca eliminar físicamente, siempre anonimizar

**Consultar:**
- `Documentacion/` (10 documentos con toda la info)
- `README.md` (setup instructions)
- `.paleta-visual-aprobada.md` (guía rápida de colores)

---

## 🚀 Próximos Pasos

### Sprint Actual: Sistema de Pagos Completo

1. **Backend:**
   - Integración completa de Stripe
   - Integración de PayPal
   - Webhooks con verificación de firma
   - Tests de concurrencia

2. **Frontend:**
   - Checkout flow completo
   - Integración Stripe Elements
   - Manejo de estados de pago
   - Recovery de pagos fallidos

Ver: `Documentacion/roadmap.md` para detalles completos

---

## 🔄 Actualizaciones de este Archivo

Cuando agregues features importantes:
1. Actualizar sección de Entidades (si hay nuevas)
2. Actualizar Endpoints Críticos
3. Actualizar Flujos Críticos
4. Mantener sincronizado con documentación principal

---

## 📝 Resumen Ejecutivo

**Migración Docker → Local (2025-11-13):**
- ✅ PostgreSQL 16 instalado y configurado localmente
- ✅ Redis 7 instalado y configurado localmente
- ✅ Backend Go como servicio systemd (sorteos-api.service)
- ✅ Docker completamente eliminado (464 MB liberados)
- ✅ Rebuild frontend: 3+ min → 10 segundos
- ✅ Debugging simplificado con journalctl
- ✅ Stack nativo, rápido y mantenible

**Stack Actual:**
- Backend: Go 1.22 (binario nativo)
- Frontend: React 18 + Vite (servido por Go)
- DB: PostgreSQL 16 (systemd)
- Cache: Redis 7 (systemd)
- Proxy: Nginx + SSL

**URLs Activas:**
- Frontend: https://sorteos.club
- API: https://sorteos.club/api/v1/
- Health: https://sorteos.club/health

---

**Última actualización:** 2025-11-13 06:45 UTC
**Versión:** 2.0 - Migración a instalación local completada
**Contacto:** Ing. Alonso Alpízar
**Despliegue:** https://sorteos.club
