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
- **Base de Datos:** PostgreSQL 15+ (transaccional)
- **Cache/Locks:** Redis 7+ (locks distribuidos, rate limiting)
- **Pagos:** Stripe (MVP) → PayPal (Fase 2)

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
│   └── migrations/           # SQL migrations
├── frontend/
│   ├── src/
│   │   ├── app/              # Router, providers
│   │   ├── features/         # auth, raffles, checkout
│   │   ├── components/ui/    # shadcn/ui components
│   │   └── lib/              # Utilidades
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

## 🛠️ Comandos Útiles

### Backend
```bash
cd backend
make run              # Ejecutar API
make test             # Tests
make migrate-up       # Aplicar migraciones
make migrate-down     # Revertir última migración
```

### Frontend
```bash
cd frontend
npm run dev           # Servidor desarrollo
npm run build         # Build producción
npm run test          # Tests (Vitest)
```

### Docker
```bash
docker-compose up -d  # Levantar todos los servicios
docker-compose logs -f api  # Ver logs
docker-compose down   # Detener
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

## ✅ Estado Actual del Sistema (2025-11-10)

### Sprint 1-2: Infraestructura y Autenticación ✅ COMPLETADO

**Despliegue:** http://62.171.188.255

#### Backend (100% ✅)
- ✅ Go 1.22 con estructura hexagonal implementada
- ✅ PostgreSQL 15 configurado y corriendo (puerto 5432)
- ✅ Redis 7 configurado y corriendo (puerto 6379)
- ✅ 3 migraciones ejecutadas:
  - `001_create_users_table` - Users con ENUMs (role, kyc_level, status)
  - `002_create_user_consents_table` - GDPR compliance
  - `003_create_audit_logs_table` - Auditoría completa
- ✅ Sistema de autenticación completo:
  - JWT (Access 15min, Refresh 7 días) con Redis
  - Bcrypt cost 12 para passwords
  - Rate limiting con Redis sliding window
  - Email verification con SendGrid
  - Audit logging en todas las acciones
- ✅ Endpoints funcionando:
  - `POST /api/v1/auth/register`
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/verify-email`
  - `POST /api/v1/auth/refresh`
  - `GET /health` - Health check
  - `GET /api/v1/ping` - Ping test

#### Frontend (100% ✅)
- ✅ React 18 + TypeScript + Vite configurado
- ✅ Tailwind CSS + shadcn/ui con **COLORES APROBADOS**
- ✅ TanStack Query + Zustand implementados
- ✅ 6 componentes UI: Button, Input, Label, Card, Alert, Badge
- ✅ 4 páginas funcionales:
  - `/login` - Login con validación Zod
  - `/register` - Registro con GDPR checkboxes
  - `/verify-email` - Verificación con código 6 dígitos
  - `/dashboard` - Dashboard protegido
- ✅ Protected routes con ProtectedRoute component
- ✅ API client con refresh automático de tokens
- ✅ Dark mode support
- ✅ Build de producción servido por Nginx

#### Infraestructura (100% ✅)
- ✅ Docker Compose configurado (postgres + redis + api)
- ✅ Nginx como reverse proxy
  - Frontend servido desde `/opt/Sorteos/frontend/dist`
  - API proxy a `localhost:8080`
  - Compresión gzip
  - Headers de seguridad
  - Cache de assets (1 año)
- ✅ Backend compilado y corriendo en Docker
- ✅ Todos los servicios healthy

#### Archivos Creados (53 total)
- **Backend:** 22 archivos (domain, use cases, repos, handlers, middlewares)
- **Frontend:** 31 archivos (components, pages, hooks, stores, config)

### 🔍 Validaciones Realizadas

```bash
# ✅ Services health
docker compose ps
# - postgres: Up 4 minutes (healthy)
# - redis: Up 4 minutes (healthy)
# - api: Up 9 seconds (healthy)

# ✅ Backend API
curl http://localhost:8080/health
# {"status":"ok","time":"2025-11-10T06:05:12Z"}

curl http://localhost:8080/api/v1/ping
# {"message":"pong","timestamp":"2025-11-10T06:05:30Z"}

# ✅ Public access
curl http://62.171.188.255/api/v1/ping
# {"message":"pong","timestamp":"2025-11-10T06:06:10Z"}

curl -I http://62.171.188.255/
# HTTP/1.1 200 OK (Frontend servido correctamente)
```

### 🔗 URLs Activas

- **Frontend**: http://62.171.188.255
- **API**: http://62.171.188.255/api/v1/
- **Health**: http://62.171.188.255/health
- **Database**: PostgreSQL en puerto 5432
- **Redis**: En puerto 6379

### 📊 Logs del Backend

```log
[INFO] Starting Sorteos Platform API (environment: development, port: 8080)
[INFO] Connected to PostgreSQL (host: postgres, database: sorteos_db)
[INFO] Connected to Redis (host: redis, db: 0)
[GIN-debug] POST /api/v1/auth/register
[GIN-debug] POST /api/v1/auth/login
[GIN-debug] POST /api/v1/auth/refresh
[GIN-debug] POST /api/v1/auth/verify-email
[INFO] Server listening (address: :8080)
```

---

## 🚀 Próximos Pasos (Sprint 3-4)

### Gestión de Sorteos (CRUD Básico)

1. **Backend:**
   - Migración `004_create_raffles_table`
   - Migración `005_create_raffle_numbers_table`
   - Domain: Raffle, RaffleNumber entities
   - Use Cases: CreateRaffle, ListRaffles, PublishRaffle
   - Implementar locks distribuidos con Redis (preparación para reservas)

2. **Frontend:**
   - Páginas: CreateRaffle, ListRaffles, RaffleDetail
   - Componentes: RaffleCard, NumberGrid
   - Form de creación con validaciones

Ver: `Documentacion/roadmap.md` para detalles completos

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

## 🔄 Actualizaciones de este Archivo

Cuando agregues features importantes:
1. Actualizar sección de Entidades (si hay nuevas)
2. Actualizar Endpoints Críticos
3. Actualizar Flujos Críticos
4. Mantener sincronizado con documentación principal

---

## 📝 Resumen Ejecutivo

**Sprint 1-2 COMPLETADO (2025-11-10):**
- ✅ 53 archivos creados (22 backend + 31 frontend)
- ✅ Sistema de autenticación funcional end-to-end
- ✅ Infraestructura desplegada y validada
- ✅ Frontend público en http://62.171.188.255
- ✅ API funcionando con rate limiting y JWT
- ✅ Base de datos con 3 migraciones aplicadas
- ✅ COLORES APROBADOS implementados (Blue #3B82F6 / Slate #64748B)

**Próximo Sprint:** Gestión de Sorteos (CRUD) + Sistema de Reservas con locks distribuidos

---

**Última actualización:** 2025-11-10 21:30 UTC
**Versión:** 1.1 - Sistema en producción
**Contacto:** Ing. Alonso Alpízar
**Despliegue:** http://62.171.188.255
