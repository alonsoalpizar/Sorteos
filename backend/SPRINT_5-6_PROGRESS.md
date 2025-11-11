# Sprint 5-6: Reservas y Pagos - Progress Report

## ✅ COMPLETADO (Backend Infrastructure - 75%)

### 1. Database Migrations ✅
- ✅ `000006_create_reservations.up/down.sql` - Tabla de reservas con TTL de 5 minutos
- ✅ `000007_create_payments.up/down.sql` - Tabla de pagos con integración Stripe
- ✅ `000008_create_idempotency_keys.up/down.sql` - Protección contra duplicados

### 2. Domain Entities ✅
- ✅ `internal/domain/entities/reservation.go` - Lógica de reservas temporales
  - NewReservation, IsExpired, CanBePaid, Confirm, Expire, Cancel
  - 5 minutos de expiración automática
- ✅ `internal/domain/entities/payment.go` - Ciclo de vida de pagos
  - NewPayment, MarkAsSucceeded, MarkAsFailed, Cancel, Refund
  - Metadata JSONB para información adicional
- ✅ `internal/domain/entities/idempotency_key.go` - Deduplicación de requests
  - NewIdempotencyKey, VerifyRequestMatch, MarkAsCompleted

### 3. Repository Layer ✅
- ✅ `internal/domain/repositories/reservation_repository.go` - Interface
- ✅ `internal/infrastructure/database/postgres_reservation_repository.go`
  - CountActiveReservationsForNumbers (detección de conflictos)
  - FindExpiredPending (para cron job)
- ✅ `internal/domain/repositories/payment_repository.go` - Interface
- ✅ `internal/infrastructure/database/postgres_payment_repository.go`
  - FindByStripePaymentIntentID (para webhooks)
- ✅ `internal/domain/repositories/idempotency_key_repository.go` - Interface
- ✅ `internal/infrastructure/database/postgres_idempotency_key_repository.go`

### 4. Distributed Locking (Redis) ✅
- ✅ `internal/infrastructure/redis/lock_service.go`
  - AcquireLock, AcquireMultipleLocks (atómico all-or-nothing)
  - Release, Extend
  - Locks con TTL automático

### 5. Payment Provider (Stripe) ✅
- ✅ `internal/infrastructure/payment/payment_provider.go` - Interface abstracta
- ✅ `internal/infrastructure/payment/stripe_provider.go` - Implementación Stripe
  - CreatePaymentIntent, GetPaymentIntent
  - ConfirmPaymentIntent, CancelPaymentIntent
  - ConstructWebhookEvent (verificación de firma)
- ✅ Stripe SDK v76 agregado a `go.mod`

### 6. Use Cases (Business Logic) ✅
- ✅ `internal/usecases/reservation_usecases.go`
  - **CreateReservation**: Locks distribuidos + validación de disponibilidad
  - **ExpireReservations**: Para cron job (libera reservas expiradas)
  - **ConfirmReservation**: Al confirmar pago
  - **CancelReservation**: Al fallar pago o cancelar usuario
  - **GetReservation**, **GetUserReservations**

- ✅ `internal/usecases/payment_usecases.go`
  - **CreatePaymentIntent**: Crea Payment Intent en Stripe + registro en DB
  - **ProcessPaymentWebhook**: Maneja eventos de Stripe
    - payment_intent.succeeded → Confirma reserva
    - payment_intent.payment_failed → Mantiene reserva pendiente
    - payment_intent.canceled → Cancela reserva
  - Soporte para Idempotency-Key

## 🔄 PENDIENTE (25% restante)

### 7. HTTP Handlers & Routes 🚧
**Archivos creados pero necesitan integración con DI:**
- `internal/adapters/http/handler/reservation/` (a crear)
- `internal/adapters/http/handler/payment/` (a crear)
- `internal/adapters/http/handler/webhook/` (a crear)

**Endpoints requeridos:**
```
POST   /api/v1/reservations          - Crear reserva
GET    /api/v1/reservations/:id      - Ver reserva
GET    /api/v1/reservations/me       - Mis reservas
POST   /api/v1/reservations/:id/cancel - Cancelar reserva

POST   /api/v1/payments/intent       - Crear payment intent
GET    /api/v1/payments/:id          - Ver pago
GET    /api/v1/payments/me           - Mis pagos

POST   /api/v1/webhooks/stripe       - Webhook de Stripe (sin auth)
```

### 8. Dependency Injection 🚧
**Necesita actualizar:**
- `cmd/api/main.go` o archivo de inicialización
- Crear instancias de:
  - ReservationRepository
  - PaymentRepository
  - IdempotencyKeyRepository
  - LockService (Redis client)
  - StripeProvider (con API key desde config)
  - ReservationUseCases
  - PaymentUseCases
- Registrar rutas con handlers

### 9. Configuration 🚧
**Agregar a `config/config.yaml` o `.env`:**
```yaml
stripe:
  secret_key: "sk_test_..."
  webhook_secret: "whsec_..."

reservations:
  expiration_minutes: 5

redis:
  host: "redis"
  port: 6379
  db: 0
```

### 10. Cron Job para Expirar Reservas ⏳
**Crear:**
- `internal/jobs/expire_reservations_job.go`
- Ejecutar cada 1 minuto
- Llamar a `reservationUseCases.ExpireReservations(ctx)`

**Opciones de implementación:**
1. Usar `github.com/robfig/cron` (recomendado)
2. Usar goroutine con ticker
3. Usar supervisor externo (cron del SO)

### 11. Frontend (50% del trabajo restante) ⏳

#### A. Number Grid Modifications
**Archivo:** `frontend/src/features/raffles/components/NumberGrid.tsx`
- ✅ Ya renderiza números disponibles/vendidos
- ⏳ Agregar selección múltiple (click para toggle)
- ⏳ Estado local de números seleccionados
- ⏳ Callback para actualizar carrito

#### B. Shopping Cart State (Zustand)
**Archivo:** `frontend/src/store/cartStore.ts` (a crear)
```typescript
interface CartStore {
  raffleId: string | null
  selectedNumbers: string[]
  addNumbers: (raffleId: string, numbers: string[]) => void
  removeNumber: (number: string) => void
  clear: () => void
  totalAmount: number
}
```

#### C. Checkout Page
**Archivo:** `frontend/src/features/checkout/pages/CheckoutPage.tsx` (a crear)
- Resumen de números seleccionados
- Desglose de precio
- Integración con Stripe Elements
- Countdown timer (5 minutos)
- Crear reserva al cargar página
- Procesar pago con Payment Intent

#### D. Payment Confirmation
**Archivo:** `frontend/src/features/checkout/pages/PaymentSuccessPage.tsx` (a crear)
- Mostrar números comprados
- Detalles del pago
- Animación de éxito (confetti)
- Botón para ver sorteo

### 12. Testing ⏳
**Concurrency Tests:**
- `internal/usecases/reservation_usecases_test.go`
- Simular 500 peticiones concurrentes
- Verificar 0 double-sales
- Verificar locks funcionan correctamente

**Integration Tests:**
- End-to-end payment flow
- Webhook processing
- Expiration job

## 📊 Arquitectura Implementada

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Layer (Pending)                     │
│  /reservations, /payments, /webhooks/stripe                 │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│                    Use Cases Layer ✅                        │
│  ReservationUseCases, PaymentUseCases                       │
│  - CreateReservation (with distributed locks)               │
│  - CreatePaymentIntent (with idempotency)                   │
│  - ProcessPaymentWebhook                                    │
│  - ExpireReservations                                       │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│                   Domain Layer ✅                            │
│  Entities: Reservation, Payment, IdempotencyKey             │
│  Repositories: Interfaces                                   │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│                Infrastructure Layer ✅                       │
│  - PostgreSQL Repositories                                  │
│  - Redis Lock Service                                       │
│  - Stripe Payment Provider                                  │
└─────────────────────────────────────────────────────────────┘
```

## 🔐 Security Features Implemented

1. **Distributed Locks** - Previene race conditions en reservas concurrentes
2. **Idempotency Keys** - Previene double-charges en pagos
3. **Webhook Signature Verification** - Valida eventos de Stripe
4. **Reservation TTL** - Libera números automáticamente (5 min)
5. **Authorization** - Verifica que user_id en token coincida con reserva/pago

## 🚀 Next Immediate Steps

1. **Integrar HTTP Handlers** (2-3 horas)
   - Crear handlers siguiendo patrón del proyecto
   - Configurar dependency injection
   - Registrar rutas en router

2. **Agregar Configuración** (30 min)
   - Stripe API keys
   - Webhook secret
   - Redis connection

3. **Implementar Cron Job** (1 hora)
   - Job de expiración de reservas
   - Logger para monitoreo

4. **Testing Backend** (2 horas)
   - Concurrency tests
   - Integration tests

5. **Frontend Checkout Flow** (6-8 horas)
   - Number selection
   - Cart state
   - Checkout page con Stripe
   - Success/error screens

6. **Deploy & Test** (2 horas)
   - Migrations en producción
   - Configurar Stripe webhook URL
   - Pruebas end-to-end

## 📝 Notas Importantes

- **Migrations**: Deben ejecutarse en producción ANTES de deploy del código
- **Stripe Test Mode**: Usar `sk_test_...` keys durante desarrollo
- **Webhook URL**: Configurar en Stripe Dashboard: `https://sorteos.club/api/v1/webhooks/stripe`
- **Redis**: Ya está configurado en docker-compose.yml
- **Lock TTL**: Debe ser >= Reservation TTL (actualmente 5 minutos para ambos)

## 🔍 Files Created in This Sprint

### Backend (15 archivos)
1. `internal/infrastructure/database/migrations/000006_create_reservations.up.sql`
2. `internal/infrastructure/database/migrations/000006_create_reservations.down.sql`
3. `internal/infrastructure/database/migrations/000007_create_payments.up.sql`
4. `internal/infrastructure/database/migrations/000007_create_payments.down.sql`
5. `internal/infrastructure/database/migrations/000008_create_idempotency_keys.up.sql`
6. `internal/infrastructure/database/migrations/000008_create_idempotency_keys.down.sql`
7. `internal/domain/entities/reservation.go`
8. `internal/domain/entities/payment.go`
9. `internal/domain/entities/idempotency_key.go`
10. `internal/domain/repositories/reservation_repository.go`
11. `internal/domain/repositories/payment_repository.go`
12. `internal/domain/repositories/idempotency_key_repository.go`
13. `internal/infrastructure/database/postgres_reservation_repository.go`
14. `internal/infrastructure/database/postgres_payment_repository.go`
15. `internal/infrastructure/database/postgres_idempotency_key_repository.go`
16. `internal/infrastructure/redis/lock_service.go`
17. `internal/infrastructure/payment/payment_provider.go`
18. `internal/infrastructure/payment/stripe_provider.go`
19. `internal/usecases/reservation_usecases.go`
20. `internal/usecases/payment_usecases.go`

### Modified
- `go.mod` - Agregado Stripe SDK v76 y lib/pq

### Pending (Frontend - 6 archivos estimados)
- `frontend/src/store/cartStore.ts`
- `frontend/src/features/checkout/pages/CheckoutPage.tsx`
- `frontend/src/features/checkout/pages/PaymentSuccessPage.tsx`
- `frontend/src/features/checkout/components/PaymentForm.tsx`
- `frontend/src/features/checkout/components/CountdownTimer.tsx`
- `frontend/src/features/raffles/components/NumberGrid.tsx` (modificar)

## 💡 Recommendations

1. **Prioridad 1**: Completar integración HTTP + DI para poder testear backend
2. **Prioridad 2**: Implementar cron job de expiración
3. **Prioridad 3**: Frontend checkout flow
4. **Prioridad 4**: Testing completo + deploy

**Estimado de tiempo restante:** 12-15 horas de desarrollo
