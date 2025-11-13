# FASE: Sistema de Reservas con Concurrencia y WebSocket

**Proyecto:** Plataforma de Sorteos  
**Fecha:** Noviembre 2025  
**Responsable:** Equipo Backend/Frontend  
**Duración estimada:** 2-3 semanas  

---

## 🎯 OBJETIVO

Implementar un sistema robusto de reserva de números de sorteo que garantice:

1. **Cero doble ventas** mediante locks distribuidos Redis + PostgreSQL
2. **Actualizaciones en tiempo real** vía WebSocket para todos los usuarios conectados
3. **Maximización de conversiones** con flujo de doble timeout (10 min selección + 5 min checkout)
4. **Sistema de pagos modular** con proveedores habilitables/deshabilitables

---

## 📊 RESULTADO ESPERADO

### Flujo completo funcional:

```
Usuario entra a sorteo
    ↓
Selecciona números (máx 10 por sesión)
    ├─ Lock inmediato en Redis + PostgreSQL
    ├─ Timer inicia: 10 minutos
    ├─ WebSocket notifica a todos: "número X reservado"
    └─ Botones "PAGAR AHORA" y "Cancelar" siempre visibles
    ↓
Usuario click "PAGAR AHORA" (en cualquier momento)
    ├─ Transición a fase checkout
    ├─ Timer se extiende: +5 minutos adicionales
    └─ Redirect a página de pago
    ↓
Si timer de selección expira (10 min):
    └─ Auto-redirect a checkout con alerta urgente
    ↓
Usuario completa pago
    ├─ Webhook confirma pago
    ├─ Números marcados como SOLD
    └─ WebSocket notifica: "números X, Y, Z vendidos"
    ↓
Si timer de checkout expira (5 min):
    ├─ Reserva marcada como EXPIRED
    ├─ Números liberados (status → AVAILABLE)
    └─ WebSocket notifica: "números liberados"
```

### Métricas de éxito:

- ✅ 0% de doble ventas en pruebas de concurrencia (1000 usuarios simultáneos)
- ✅ Latencia de actualización WebSocket < 100ms
- ✅ Tasa de conversión reserva → pago > 70%
- ✅ Tiempo de respuesta API < 200ms (P95)

---

## 🏗️ ARQUITECTURA

### Componentes principales:

```
┌─────────────────────────────────────────────────────────┐
│                      FRONTEND                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ NumberGrid   │  │ WebSocket    │  │ Reservation  │  │
│  │ Component    │←─┤ Client       │  │ Timer        │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
                            ↕ WebSocket
┌─────────────────────────────────────────────────────────┐
│                      BACKEND                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Reservation  │  │ WebSocket    │  │ Payment      │  │
│  │ UseCases     │→─┤ Hub          │  │ Providers    │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
│         ↓                                     ↓          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Redis Locks  │  │ PostgreSQL   │  │ BAC/SINPE    │  │
│  │ (30s TTL)    │  │ Transactions │  │ Adapters     │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
```

---

## 📦 TAREAS DETALLADAS

## SEMANA 1: Core del Sistema de Reservas

### TAREA 1.1: Modelo de Datos - Reservas con Doble Timeout

**Archivo:** `backend/internal/domain/entities/reservation.go`

**Objetivo:** Crear entidad Reservation con soporte para dos fases (selection + checkout)

**Implementación:**

```go
package entities

import (
    "errors"
    "time"
    "github.com/google/uuid"
)

const (
    MaxNumbersPerReservation = 10
    ReservationSelectionTimeout = 10 * time.Minute
    ReservationCheckoutTimeout  = 5 * time.Minute
)

type ReservationPhase string

const (
    ReservationPhaseSelection ReservationPhase = "selection"
    ReservationPhaseCheckout  ReservationPhase = "checkout"
    ReservationPhaseCompleted ReservationPhase = "completed"
    ReservationPhaseExpired   ReservationPhase = "expired"
)

type ReservationStatus string

const (
    ReservationStatusPending   ReservationStatus = "pending"
    ReservationStatusConfirmed ReservationStatus = "confirmed"
    ReservationStatusExpired   ReservationStatus = "expired"
    ReservationStatusCancelled ReservationStatus = "cancelled"
)

type Reservation struct {
    ID                  uuid.UUID
    RaffleID            uuid.UUID
    UserID              uuid.UUID
    Numbers             []string
    SessionID           string
    TotalAmount         float64
    
    // Fases y timeouts
    Phase               ReservationPhase
    Status              ReservationStatus
    SelectionStartedAt  time.Time
    CheckoutStartedAt   *time.Time
    ExpiresAt           time.Time
    
    // Tracking
    CreatedAt           time.Time
    UpdatedAt           time.Time
    CompletedAt         *time.Time
}

// NewReservation crea una nueva reserva en fase de selección
func NewReservation(raffleID, userID uuid.UUID, numbers []string, sessionID string, totalAmount float64) (*Reservation, error) {
    if len(numbers) == 0 {
        return nil, errors.New("must reserve at least one number")
    }
    
    if len(numbers) > MaxNumbersPerReservation {
        return nil, errors.New("cannot reserve more than 10 numbers per session")
    }
    
    now := time.Now()
    
    return &Reservation{
        ID:                 uuid.New(),
        RaffleID:           raffleID,
        UserID:             userID,
        Numbers:            numbers,
        SessionID:          sessionID,
        TotalAmount:        totalAmount,
        Phase:              ReservationPhaseSelection,
        Status:             ReservationStatusPending,
        SelectionStartedAt: now,
        ExpiresAt:          now.Add(ReservationSelectionTimeout),
        CreatedAt:          now,
        UpdatedAt:          now,
    }, nil
}

// AddNumber agrega un número a la reserva existente
func (r *Reservation) AddNumber(numberID string) error {
    if r.Phase != ReservationPhaseSelection {
        return errors.New("can only add numbers during selection phase")
    }
    
    if len(r.Numbers) >= MaxNumbersPerReservation {
        return errors.New("maximum 10 numbers per reservation")
    }
    
    // Verificar duplicados
    for _, num := range r.Numbers {
        if num == numberID {
            return errors.New("number already in reservation")
        }
    }
    
    r.Numbers = append(r.Numbers, numberID)
    r.UpdatedAt = time.Now()
    
    return nil
}

// MoveToCheckout transiciona la reserva a fase de checkout
func (r *Reservation) MoveToCheckout() error {
    if r.Phase != ReservationPhaseSelection {
        return errors.New("reservation not in selection phase")
    }
    
    if r.IsExpired() {
        return errors.New("reservation has expired")
    }
    
    now := time.Now()
    r.Phase = ReservationPhaseCheckout
    r.CheckoutStartedAt = &now
    r.ExpiresAt = now.Add(ReservationCheckoutTimeout)
    r.UpdatedAt = now
    
    return nil
}

// Confirm marca la reserva como completada
func (r *Reservation) Confirm() error {
    if r.Status != ReservationStatusPending {
        return errors.New("can only confirm pending reservations")
    }
    
    now := time.Now()
    r.Phase = ReservationPhaseCompleted
    r.Status = ReservationStatusConfirmed
    r.CompletedAt = &now
    r.UpdatedAt = now
    
    return nil
}

// Cancel cancela la reserva
func (r *Reservation) Cancel() error {
    if r.Status != ReservationStatusPending {
        return errors.New("can only cancel pending reservations")
    }
    
    r.Status = ReservationStatusCancelled
    r.UpdatedAt = time.Now()
    
    return nil
}

// Expire marca la reserva como expirada
func (r *Reservation) Expire() error {
    if r.Status != ReservationStatusPending {
        return errors.New("can only expire pending reservations")
    }
    
    r.Phase = ReservationPhaseExpired
    r.Status = ReservationStatusExpired
    r.UpdatedAt = time.Now()
    
    return nil
}

// IsExpired verifica si la reserva ha expirado
func (r *Reservation) IsExpired() bool {
    return time.Now().After(r.ExpiresAt)
}

// TimeRemaining retorna el tiempo restante antes de expiración
func (r *Reservation) TimeRemaining() time.Duration {
    remaining := time.Until(r.ExpiresAt)
    if remaining < 0 {
        return 0
    }
    return remaining
}
```

**Migración SQL:**

```sql
-- backend/migrations/006_enhance_reservations_table.up.sql

-- Modificar tabla existente
ALTER TABLE reservations 
ADD COLUMN phase VARCHAR(20) DEFAULT 'selection',
ADD COLUMN selection_started_at TIMESTAMP DEFAULT NOW(),
ADD COLUMN checkout_started_at TIMESTAMP NULL,
ADD CONSTRAINT check_phase CHECK (phase IN ('selection', 'checkout', 'completed', 'expired'));

-- Índices para queries frecuentes
CREATE INDEX idx_reservations_phase ON reservations(phase);
CREATE INDEX idx_reservations_expires_at ON reservations(expires_at) WHERE status = 'pending';
CREATE INDEX idx_reservations_session_id ON reservations(session_id);

-- Down migration
-- backend/migrations/006_enhance_reservations_table.down.sql
ALTER TABLE reservations 
DROP COLUMN phase,
DROP COLUMN selection_started_at,
DROP COLUMN checkout_started_at;

DROP INDEX IF EXISTS idx_reservations_phase;
DROP INDEX IF EXISTS idx_reservations_expires_at;
DROP INDEX IF EXISTS idx_reservations_session_id;
```

**Testing:**

```go
// backend/internal/domain/entities/reservation_test.go

func TestReservation_AddNumber(t *testing.T) {
    r, _ := NewReservation(uuid.New(), uuid.New(), []string{"00"}, "session-1", 100.0)
    
    // Agregar números válidos
    err := r.AddNumber("01")
    assert.NoError(t, err)
    assert.Len(t, r.Numbers, 2)
    
    // Intentar agregar más de 10
    for i := 2; i < 12; i++ {
        r.AddNumber(fmt.Sprintf("%02d", i))
    }
    err = r.AddNumber("12")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "maximum 10 numbers")
    
    // Intentar agregar duplicado
    err = r.AddNumber("00")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already in reservation")
}

func TestReservation_MoveToCheckout(t *testing.T) {
    r, _ := NewReservation(uuid.New(), uuid.New(), []string{"00"}, "session-1", 100.0)
    
    // Transición exitosa
    err := r.MoveToCheckout()
    assert.NoError(t, err)
    assert.Equal(t, ReservationPhaseCheckout, r.Phase)
    assert.NotNil(t, r.CheckoutStartedAt)
    
    // Verificar que expiración se extendió
    expectedExpiry := r.CheckoutStartedAt.Add(ReservationCheckoutTimeout)
    assert.WithinDuration(t, expectedExpiry, r.ExpiresAt, time.Second)
    
    // No se puede hacer checkout dos veces
    err = r.MoveToCheckout()
    assert.Error(t, err)
}

func TestReservation_IsExpired(t *testing.T) {
    r, _ := NewReservation(uuid.New(), uuid.New(), []string{"00"}, "session-1", 100.0)
    
    // No expirada al inicio
    assert.False(t, r.IsExpired())
    
    // Forzar expiración
    r.ExpiresAt = time.Now().Add(-1 * time.Minute)
    assert.True(t, r.IsExpired())
}
```

**Resultado esperado:**
- ✅ Entidad Reservation con dos fases implementada
- ✅ Tests unitarios pasando al 100%
- ✅ Migración SQL aplicada correctamente
- ✅ Límite de 10 números validado

---

### TAREA 1.2: Sistema de Locks Distribuidos con Redis

**Archivo:** `backend/internal/infrastructure/redis/lock_service.go`

**Objetivo:** Implementar locks distribuidos para evitar doble venta de números

**Implementación:**

```go
package redis

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
)

var (
    ErrLockNotAcquired = errors.New("could not acquire lock")
    ErrLockNotReleased = errors.New("could not release lock")
)

type LockService struct {
    client *redis.Client
}

func NewLockService(client *redis.Client) *LockService {
    return &LockService{client: client}
}

// AcquireMultipleLocks intenta adquirir locks para múltiples números
// Retorna error si ALGUNO de los locks no se puede adquirir
func (s *LockService) AcquireMultipleLocks(ctx context.Context, keys []string, ttl time.Duration) ([]*Lock, error) {
    locks := make([]*Lock, 0, len(keys))
    acquiredKeys := make([]string, 0, len(keys))
    
    // Intentar adquirir todos los locks
    for _, key := range keys {
        lock, err := s.AcquireLock(ctx, key, ttl)
        if err != nil {
            // Si falla uno, liberar todos los ya adquiridos
            s.releaseKeys(ctx, acquiredKeys)
            return nil, fmt.Errorf("%w: key=%s", ErrLockNotAcquired, key)
        }
        
        locks = append(locks, lock)
        acquiredKeys = append(acquiredKeys, key)
    }
    
    return locks, nil
}

// AcquireLock intenta adquirir un lock individual
func (s *LockService) AcquireLock(ctx context.Context, key string, ttl time.Duration) (*Lock, error) {
    value := fmt.Sprintf("lock:%d", time.Now().UnixNano())
    
    // SETNX (SET if Not eXists) con TTL
    ok, err := s.client.SetNX(ctx, key, value, ttl).Result()
    if err != nil {
        return nil, fmt.Errorf("redis error: %w", err)
    }
    
    if !ok {
        return nil, ErrLockNotAcquired
    }
    
    return &Lock{
        Key:      key,
        Value:    value,
        ExpireAt: time.Now().Add(ttl),
    }, nil
}

// ReleaseLock libera un lock específico
func (s *LockService) ReleaseLock(ctx context.Context, lock *Lock) error {
    // Solo liberar si el valor coincide (evitar liberar locks de otros)
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
    
    result, err := s.client.Eval(ctx, script, []string{lock.Key}, lock.Value).Result()
    if err != nil {
        return fmt.Errorf("redis error: %w", err)
    }
    
    if result.(int64) == 0 {
        return ErrLockNotReleased
    }
    
    return nil
}

// ReleaseMultipleLocks libera múltiples locks
func ReleaseMultipleLocks(ctx context.Context, locks []*Lock) error {
    // Nota: No retornar error si alguno falla (los locks tienen TTL)
    for _, lock := range locks {
        // Ignorar errores individuales
        _ = lock.Service.ReleaseLock(ctx, lock)
    }
    return nil
}

// releaseKeys es un helper interno para liberar keys en caso de fallo
func (s *LockService) releaseKeys(ctx context.Context, keys []string) {
    for _, key := range keys {
        s.client.Del(ctx, key)
    }
}

// ReservationLockKey genera la key de Redis para un número de sorteo
func ReservationLockKey(raffleID, numberID string) string {
    return fmt.Sprintf("lock:raffle:%s:number:%s", raffleID, numberID)
}

type Lock struct {
    Key      string
    Value    string
    ExpireAt time.Time
    Service  *LockService
}
```

**Testing:**

```go
// backend/internal/infrastructure/redis/lock_service_test.go

func TestLockService_AcquireMultipleLocks_Success(t *testing.T) {
    ctx := context.Background()
    client := setupTestRedis(t)
    service := NewLockService(client)
    
    keys := []string{"lock:test:1", "lock:test:2", "lock:test:3"}
    
    locks, err := service.AcquireMultipleLocks(ctx, keys, 30*time.Second)
    
    assert.NoError(t, err)
    assert.Len(t, locks, 3)
    
    // Verificar que los locks están en Redis
    for _, key := range keys {
        val, err := client.Get(ctx, key).Result()
        assert.NoError(t, err)
        assert.NotEmpty(t, val)
    }
    
    // Cleanup
    ReleaseMultipleLocks(ctx, locks)
}

func TestLockService_AcquireMultipleLocks_Conflict(t *testing.T) {
    ctx := context.Background()
    client := setupTestRedis(t)
    service := NewLockService(client)
    
    // Usuario A adquiere lock del número 1
    client.Set(ctx, "lock:test:1", "existing", 30*time.Second)
    
    // Usuario B intenta adquirir números 1, 2, 3
    keys := []string{"lock:test:1", "lock:test:2", "lock:test:3"}
    
    locks, err := service.AcquireMultipleLocks(ctx, keys, 30*time.Second)
    
    // Debe fallar
    assert.Error(t, err)
    assert.Nil(t, locks)
    assert.ErrorIs(t, err, ErrLockNotAcquired)
    
    // Verificar que no se adquirieron locks parciales
    val, err := client.Get(ctx, "lock:test:2").Result()
    assert.Error(t, err) // No debe existir
    assert.Empty(t, val)
}

func TestLockService_ConcurrentAcquisition(t *testing.T) {
    ctx := context.Background()
    client := setupTestRedis(t)
    service := NewLockService(client)
    
    key := "lock:test:concurrent"
    
    successCount := atomic.Int32{}
    wg := sync.WaitGroup{}
    
    // 100 goroutines intentan adquirir el mismo lock
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            lock, err := service.AcquireLock(ctx, key, 30*time.Second)
            if err == nil {
                successCount.Add(1)
                service.ReleaseLock(ctx, lock)
            }
        }()
    }
    
    wg.Wait()
    
    // Solo uno debe tener éxito
    assert.Equal(t, int32(1), successCount.Load())
}
```

**Resultado esperado:**
- ✅ Locks distribuidos implementados con SETNX
- ✅ All-or-nothing acquisition (falla si algún lock no se puede adquirir)
- ✅ Tests de concurrencia pasando
- ✅ TTL automático en Redis

---

### TAREA 1.3: Use Case - Crear Reserva con Locks

**Archivo:** `backend/internal/usecases/reservation_usecases.go`

**Objetivo:** Implementar lógica de negocio para crear reservas con protección de concurrencia

**Implementación actualizada:**

```go
package usecases

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    
    "github.com/sorteos-platform/backend/internal/domain/entities"
    "github.com/sorteos-platform/backend/internal/domain/repositories"
    "github.com/sorteos-platform/backend/internal/infrastructure/redis"
)

type ReservationUseCases struct {
    reservationRepo repositories.ReservationRepository
    raffleRepo      repositories.RaffleRepository
    numberRepo      repositories.RaffleNumberRepository
    lockService     *redis.LockService
    wsHub           *websocket.Hub  // Para notificaciones
}

func NewReservationUseCases(
    reservationRepo repositories.ReservationRepository,
    raffleRepo repositories.RaffleRepository,
    numberRepo repositories.RaffleNumberRepository,
    lockService *redis.LockService,
    wsHub *websocket.Hub,
) *ReservationUseCases {
    return &ReservationUseCases{
        reservationRepo: reservationRepo,
        raffleRepo:      raffleRepo,
        numberRepo:      numberRepo,
        lockService:     lockService,
        wsHub:           wsHub,
    }
}

type CreateReservationInput struct {
    RaffleID  uuid.UUID
    UserID    uuid.UUID
    NumberIDs []string
    SessionID string
}

// CreateReservation crea una nueva reserva con locks distribuidos
func (uc *ReservationUseCases) CreateReservation(ctx context.Context, input CreateReservationInput) (*entities.Reservation, error) {
    // 1. Validar límite de números
    if len(input.NumberIDs) > entities.MaxNumbersPerReservation {
        return nil, fmt.Errorf("cannot reserve more than %d numbers", entities.MaxNumbersPerReservation)
    }
    
    // 2. Verificar idempotencia (evitar duplicados por sesión)
    existingReservation, err := uc.reservationRepo.FindBySessionID(ctx, input.SessionID)
    if err != nil {
        return nil, fmt.Errorf("error checking existing reservation: %w", err)
    }
    if existingReservation != nil && !existingReservation.IsExpired() {
        // Ya existe una reserva activa para esta sesión
        return existingReservation, nil
    }
    
    // 3. Validar que el sorteo existe y está activo
    raffle, err := uc.raffleRepo.FindByID(ctx, input.RaffleID)
    if err != nil {
        return nil, fmt.Errorf("error fetching raffle: %w", err)
    }
    if raffle == nil {
        return nil, errors.New("raffle not found")
    }
    if raffle.Status != "active" {
        return nil, errors.New("raffle is not active")
    }
    
    // 4. Adquirir locks distribuidos para todos los números
    lockKeys := make([]string, len(input.NumberIDs))
    for i, numberID := range input.NumberIDs {
        lockKeys[i] = redis.ReservationLockKey(input.RaffleID.String(), numberID)
    }
    
    locks, err := uc.lockService.AcquireMultipleLocks(ctx, lockKeys, 30*time.Second)
    if err != nil {
        if errors.Is(err, redis.ErrLockNotAcquired) {
            return nil, errors.New("one or more numbers are already reserved")
        }
        return nil, fmt.Errorf("error acquiring locks: %w", err)
    }
    
    // Liberar locks al finalizar (error o éxito)
    defer redis.ReleaseMultipleLocks(ctx, locks)
    
    // 5. Verificar en DB que los números NO están reservados/vendidos
    numbers, err := uc.numberRepo.FindByIDs(ctx, input.NumberIDs)
    if err != nil {
        return nil, fmt.Errorf("error fetching numbers: %w", err)
    }
    
    for _, num := range numbers {
        if !num.IsAvailable() {
            return nil, fmt.Errorf("number %s is not available", num.Number)
        }
    }
    
    // 6. Calcular monto total
    pricePerNumber, _ := raffle.PricePerNumber.Float64()
    totalAmount := float64(len(input.NumberIDs)) * pricePerNumber
    
    // 7. Crear entidad Reservation
    reservation, err := entities.NewReservation(
        input.RaffleID,
        input.UserID,
        input.NumberIDs,
        input.SessionID,
        totalAmount,
    )
    if err != nil {
        return nil, fmt.Errorf("error creating reservation entity: %w", err)
    }
    
    // 8. Guardar en base de datos (dentro de transacción)
    err = uc.reservationRepo.CreateWithTransaction(ctx, func(txCtx context.Context) error {
        // Crear reserva
        if err := uc.reservationRepo.Create(txCtx, reservation); err != nil {
            return err
        }
        
        // Actualizar estado de números a RESERVED
        for _, numberID := range input.NumberIDs {
            if err := uc.numberRepo.MarkAsReserved(txCtx, numberID, reservation.ID, reservation.ExpiresAt); err != nil {
                return err
            }
        }
        
        return nil
    })
    
    if err != nil {
        return nil, fmt.Errorf("error saving reservation: %w", err)
    }
    
    // 9. Notificar via WebSocket a todos los usuarios conectados
    for _, numberID := range input.NumberIDs {
        uc.wsHub.BroadcastNumberUpdate(
            input.RaffleID.String(),
            numberID,
            "reserved",
            &input.UserID,
        )
    }
    
    return reservation, nil
}

// AddNumberToReservation agrega un número a una reserva existente
func (uc *ReservationUseCases) AddNumberToReservation(ctx context.Context, reservationID uuid.UUID, numberID string) error {
    // 1. Obtener reserva
    reservation, err := uc.reservationRepo.FindByID(ctx, reservationID)
    if err != nil {
        return fmt.Errorf("error fetching reservation: %w", err)
    }
    if reservation == nil {
        return errors.New("reservation not found")
    }
    
    // 2. Validar que está en fase de selección y no expirada
    if reservation.Phase != entities.ReservationPhaseSelection {
        return errors.New("can only add numbers during selection phase")
    }
    if reservation.IsExpired() {
        return errors.New("reservation has expired")
    }
    
    // 3. Adquirir lock del número
    lockKey := redis.ReservationLockKey(reservation.RaffleID.String(), numberID)
    lock, err := uc.lockService.AcquireLock(ctx, lockKey, 30*time.Second)
    if err != nil {
        if errors.Is(err, redis.ErrLockNotAcquired) {
            return errors.New("number is already reserved")
        }
        return fmt.Errorf("error acquiring lock: %w", err)
    }
    defer uc.lockService.ReleaseLock(ctx, lock)
    
    // 4. Verificar disponibilidad en DB
    number, err := uc.numberRepo.FindByID(ctx, numberID)
    if err != nil {
        return fmt.Errorf("error fetching number: %w", err)
    }
    if !number.IsAvailable() {
        return errors.New("number is not available")
    }
    
    // 5. Agregar número a la reserva
    if err := reservation.AddNumber(numberID); err != nil {
        return err
    }
    
    // 6. Actualizar en DB
    err = uc.reservationRepo.UpdateWithTransaction(ctx, func(txCtx context.Context) error {
        if err := uc.reservationRepo.Update(txCtx, reservation); err != nil {
            return err
        }
        
        if err := uc.numberRepo.MarkAsReserved(txCtx, numberID, reservation.ID, reservation.ExpiresAt); err != nil {
            return err
        }
        
        return nil
    })
    
    if err != nil {
        return fmt.Errorf("error updating reservation: %w", err)
    }
    
    // 7. Notificar via WebSocket
    uc.wsHub.BroadcastNumberUpdate(
        reservation.RaffleID.String(),
        numberID,
        "reserved",
        &reservation.UserID,
    )
    
    return nil
}

// MoveToCheckout transiciona la reserva a fase de checkout
func (uc *ReservationUseCases) MoveToCheckout(ctx context.Context, reservationID uuid.UUID) error {
    reservation, err := uc.reservationRepo.FindByID(ctx, reservationID)
    if err != nil {
        return fmt.Errorf("error fetching reservation: %w", err)
    }
    if reservation == nil {
        return errors.New("reservation not found")
    }
    
    if err := reservation.MoveToCheckout(); err != nil {
        return err
    }
    
    if err := uc.reservationRepo.Update(ctx, reservation); err != nil {
        return fmt.Errorf("error updating reservation: %w", err)
    }
    
    return nil
}

// CancelReservation cancela una reserva y libera números
func (uc *ReservationUseCases) CancelReservation(ctx context.Context, reservationID uuid.UUID) error {
    reservation, err := uc.reservationRepo.FindByID(ctx, reservationID)
    if err != nil {
        return fmt.Errorf("error fetching reservation: %w", err)
    }
    if reservation == nil {
        return errors.New("reservation not found")
    }
    
    if err := reservation.Cancel(); err != nil {
        return err
    }
    
    // Liberar números
    err = uc.reservationRepo.UpdateWithTransaction(ctx, func(txCtx context.Context) error {
        if err := uc.reservationRepo.Update(txCtx, reservation); err != nil {
            return err
        }
        
        for _, numberID := range reservation.Numbers {
            if err := uc.numberRepo.MarkAsAvailable(txCtx, numberID); err != nil {
                return err
            }
        }
        
        return nil
    })
    
    if err != nil {
        return fmt.Errorf("error cancelling reservation: %w", err)
    }
    
    // Notificar liberación via WebSocket
    for _, numberID := range reservation.Numbers {
        uc.wsHub.BroadcastNumberUpdate(
            reservation.RaffleID.String(),
            numberID,
            "available",
            nil,
        )
    }
    
    return nil
}
```

**Resultado esperado:**
- ✅ Reservas creadas con locks distribuidos
- ✅ Protección contra race conditions
- ✅ Transacciones atómicas en PostgreSQL
- ✅ Notificaciones WebSocket integradas

---

## SEMANA 2: WebSocket en Tiempo Real

### TAREA 2.1: WebSocket Hub (Backend)

**Archivo:** `backend/internal/infrastructure/websocket/hub.go`

**Objetivo:** Implementar hub central para manejar conexiones WebSocket y broadcast de mensajes

**Implementación:**

```go
package websocket

import (
    "encoding/json"
    "log"
    "sync"
)

type MessageType string

const (
    MessageTypeNumberUpdate MessageType = "number_update"
    MessageTypeReservationExpired MessageType = "reservation_expired"
)

type Message struct {
    Type     MessageType            `json:"type"`
    RaffleID string                 `json:"raffle_id"`
    Data     map[string]interface{} `json:"data"`
}

type Hub struct {
    // Clientes organizados por raffle_id
    raffles map[string]map[*Client]bool
    mu      sync.RWMutex
    
    // Channels
    broadcast  chan *Message
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        raffles:    make(map[string]map[*Client]bool),
        broadcast:  make(chan *Message, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

// Run inicia el loop principal del hub (debe correr en goroutine)
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.registerClient(client)
            
        case client := <-h.unregister:
            h.unregisterClient(client)
            
        case message := <-h.broadcast:
            h.broadcastToRaffle(message)
        }
    }
}

func (h *Hub) registerClient(client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()
    
    if h.raffles[client.raffleID] == nil {
        h.raffles[client.raffleID] = make(map[*Client]bool)
    }
    
    h.raffles[client.raffleID][client] = true
    
    log.Printf("[WebSocket] Client registered to raffle %s (total: %d)", 
        client.raffleID, len(h.raffles[client.raffleID]))
}

func (h *Hub) unregisterClient(client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()
    
    if clients, ok := h.raffles[client.raffleID]; ok {
        if _, exists := clients[client]; exists {
            delete(clients, client)
            close(client.send)
            
            // Si no quedan clientes, eliminar el raffle del mapa
            if len(clients) == 0 {
                delete(h.raffles, client.raffleID)
            }
            
            log.Printf("[WebSocket] Client unregistered from raffle %s (remaining: %d)", 
                client.raffleID, len(clients))
        }
    }
}

func (h *Hub) broadcastToRaffle(message *Message) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    clients, ok := h.raffles[message.RaffleID]
    if !ok {
        return // No hay clientes conectados a este sorteo
    }
    
    messageJSON, err := json.Marshal(message)
    if err != nil {
        log.Printf("[WebSocket] Error marshaling message: %v", err)
        return
    }
    
    // Enviar a todos los clientes del sorteo
    for client := range clients {
        select {
        case client.send <- messageJSON:
        default:
            // Canal lleno, cerrar cliente
            close(client.send)
            delete(clients, client)
        }
    }
}

// BroadcastNumberUpdate notifica cambio de estado de un número
func (h *Hub) BroadcastNumberUpdate(raffleID, numberID, status string, userID *string) {
    data := map[string]interface{}{
        "number_id": numberID,
        "status":    status,
    }
    
    if userID != nil {
        data["user_id"] = *userID
    }
    
    h.broadcast <- &Message{
        Type:     MessageTypeNumberUpdate,
        RaffleID: raffleID,
        Data:     data,
    }
}

// BroadcastReservationExpired notifica que una reserva expiró
func (h *Hub) BroadcastReservationExpired(raffleID string, numberIDs []string) {
    h.broadcast <- &Message{
        Type:     MessageTypeReservationExpired,
        RaffleID: raffleID,
        Data: map[string]interface{}{
            "number_ids": numberIDs,
        },
    }
}

// GetConnectedClients retorna el número de clientes conectados a un sorteo
func (h *Hub) GetConnectedClients(raffleID string) int {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    if clients, ok := h.raffles[raffleID]; ok {
        return len(clients)
    }
    return 0
}
```

**Archivo:** `backend/internal/infrastructure/websocket/client.go`

```go
package websocket

import (
    "log"
    "time"
    
    "github.com/gorilla/websocket"
)

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512
)

type Client struct {
    hub      *Hub
    conn     *websocket.Conn
    send     chan []byte
    raffleID string
}

// ReadPump lee mensajes del WebSocket (principalmente para mantener la conexión)
func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })
    
    for {
        _, _, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("[WebSocket] Unexpected close error: %v", err)
            }
            break
        }
    }
}

// WritePump escribe mensajes al WebSocket
func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // Hub cerró el canal
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }
            
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

**Resultado esperado:**
- ✅ Hub WebSocket funcionando con múltiples clientes
- ✅ Broadcast eficiente por raffle_id
- ✅ Ping/Pong automático para mantener conexiones
- ✅ Limpieza automática de clientes desconectados

---

### TAREA 2.2: WebSocket HTTP Handler

**Archivo:** `backend/internal/adapters/http/handlers/websocket_handler.go`

**Objetivo:** Endpoint HTTP que upgradea a WebSocket

**Implementación:**

```go
package handlers

import (
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    
    ws "github.com/sorteos-platform/backend/internal/infrastructure/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // En producción: validar origen específico
        // return r.Header.Get("Origin") == "https://sorteos.com"
        return true
    },
}

type WebSocketHandler struct {
    hub *ws.Hub
}

func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
    return &WebSocketHandler{hub: hub}
}

// HandleConnection maneja la conexión WebSocket para un sorteo específico
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
    raffleID := c.Param("id")
    
    if raffleID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "raffle_id is required"})
        return
    }
    
    // Upgrade HTTP connection a WebSocket
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("[WebSocket] Failed to upgrade connection: %v", err)
        return
    }
    
    // Crear cliente
    client := &ws.Client{
        hub:      h.hub,
        conn:     conn,
        send:     make(chan []byte, 256),
        raffleID: raffleID,
    }
    
    // Registrar cliente en el hub
    h.hub.register <- client
    
    // Iniciar pumps en goroutines separadas
    go client.WritePump()
    go client.ReadPump()
}

// GetStats retorna estadísticas de conexiones (para debug/monitoring)
func (h *WebSocketHandler) GetStats(c *gin.Context) {
    raffleID := c.Param("id")
    
    connectedClients := h.hub.GetConnectedClients(raffleID)
    
    c.JSON(http.StatusOK, gin.H{
        "raffle_id":         raffleID,
        "connected_clients": connectedClients,
    })
}
```

**Registrar rutas:**

```go
// backend/cmd/api/main.go

func setupRoutes(router *gin.Engine, deps *dependencies) {
    // ... rutas existentes
    
    // WebSocket
    wsHandler := handlers.NewWebSocketHandler(deps.wsHub)
    router.GET("/raffles/:id/ws", wsHandler.HandleConnection)
    router.GET("/raffles/:id/ws/stats", wsHandler.GetStats)
}

// Iniciar hub en main
func main() {
    // ... setup existente
    
    wsHub := websocket.NewHub()
    go wsHub.Run() // ← Importante: correr en goroutine
    
    // ... resto del setup
}
```

**Resultado esperado:**
- ✅ Endpoint `/raffles/:id/ws` funcionando
- ✅ Clientes pueden conectarse vía WebSocket
- ✅ Hub corriendo en background

---

### TAREA 2.3: WebSocket Client (Frontend)

**Archivo:** `frontend/src/hooks/useRaffleWebSocket.ts`

**Objetivo:** Hook de React para conectar y escuchar actualizaciones WebSocket

**Implementación:**

```typescript
import { useEffect, useRef, useState, useCallback } from 'react';
import { toast } from 'sonner';

interface NumberUpdate {
  number_id: string;
  status: 'available' | 'reserved' | 'sold';
  user_id?: string;
}

interface ReservationExpired {
  number_ids: string[];
}

interface WebSocketMessage {
  type: 'number_update' | 'reservation_expired';
  raffle_id: string;
  data: NumberUpdate | ReservationExpired;
}

interface UseRaffleWebSocketReturn {
  isConnected: boolean;
  onNumberUpdate: (callback: (update: NumberUpdate) => void) => void;
  onReservationExpired: (callback: (data: ReservationExpired) => void) => void;
}

export function useRaffleWebSocket(raffleId: string): UseRaffleWebSocketReturn {
  const ws = useRef<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const reconnectTimeout = useRef<NodeJS.Timeout>();
  
  const numberUpdateCallbacks = useRef<((update: NumberUpdate) => void)[]>([]);
  const reservationExpiredCallbacks = useRef<((data: ReservationExpired) => void)[]>([]);

  const connect = useCallback(() => {
    const wsUrl = `${import.meta.env.VITE_WS_URL}/raffles/${raffleId}/ws`;
    
    ws.current = new WebSocket(wsUrl);
    
    ws.current.onopen = () => {
      console.log('[WebSocket] Connected to raffle', raffleId);
      setIsConnected(true);
      toast.success('Conexión en tiempo real activada');
    };
    
    ws.current.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        
        if (message.type === 'number_update') {
          const data = message.data as NumberUpdate;
          numberUpdateCallbacks.current.forEach(cb => cb(data));
        } else if (message.type === 'reservation_expired') {
          const data = message.data as ReservationExpired;
          reservationExpiredCallbacks.current.forEach(cb => cb(data));
        }
      } catch (error) {
        console.error('[WebSocket] Error parsing message:', error);
      }
    };
    
    ws.current.onclose = () => {
      console.log('[WebSocket] Disconnected');
      setIsConnected(false);
      
      // Auto-reconnect después de 3 segundos
      reconnectTimeout.current = setTimeout(() => {
        console.log('[WebSocket] Attempting to reconnect...');
        connect();
      }, 3000);
    };
    
    ws.current.onerror = (error) => {
      console.error('[WebSocket] Error:', error);
      toast.error('Error en conexión en tiempo real');
    };
  }, [raffleId]);

  useEffect(() => {
    connect();
    
    return () => {
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }
      
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [connect]);

  const onNumberUpdate = useCallback((callback: (update: NumberUpdate) => void) => {
    numberUpdateCallbacks.current.push(callback);
  }, []);

  const onReservationExpired = useCallback((callback: (data: ReservationExpired) => void) => {
    reservationExpiredCallbacks.current.push(callback);
  }, []);

  return {
    isConnected,
    onNumberUpdate,
    onReservationExpired,
  };
}
```

**Archivo:** `frontend/src/components/RaffleNumberGrid.tsx`

**Uso del hook:**

```typescript
import { useState, useEffect } from 'react';
import { useRaffleWebSocket } from '@/hooks/useRaffleWebSocket';
import { RaffleNumber } from '@/types';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

interface Props {
  raffleId: string;
  initialNumbers: RaffleNumber[];
  onSelect: (numberId: string) => void;
  selectedNumbers: string[];
}

export function RaffleNumberGrid({ 
  raffleId, 
  initialNumbers, 
  onSelect,
  selectedNumbers 
}: Props) {
  const [numbers, setNumbers] = useState<RaffleNumber[]>(initialNumbers);
  const { isConnected, onNumberUpdate, onReservationExpired } = useRaffleWebSocket(raffleId);

  // Escuchar actualizaciones de números
  useEffect(() => {
    onNumberUpdate((update) => {
      setNumbers(prev => 
        prev.map(num => 
          num.id === update.number_id
            ? { ...num, status: update.status }
            : num
        )
      );
    });

    onReservationExpired((data) => {
      // Números liberados por expiración
      setNumbers(prev =>
        prev.map(num =>
          data.number_ids.includes(num.id)
            ? { ...num, status: 'available' }
            : num
        )
      );
      
      toast.info(`${data.number_ids.length} números han sido liberados`);
    });
  }, [onNumberUpdate, onReservationExpired]);

  return (
    <div className="space-y-4">
      {/* Indicador de conexión */}
      <div className="flex items-center gap-2 text-sm">
        <div className={cn(
          "h-2 w-2 rounded-full",
          isConnected ? "bg-success animate-pulse" : "bg-neutral-300"
        )} />
        <span className="text-neutral-600">
          {isConnected ? 'Actualizaciones en tiempo real' : 'Reconectando...'}
        </span>
      </div>

      {/* Grid de números */}
      <div className="grid grid-cols-10 gap-2">
        {numbers.map(num => (
          <Button
            key={num.id}
            variant={
              selectedNumbers.includes(num.id) ? 'default' :
              num.status === 'available' ? 'outline' : 'secondary'
            }
            disabled={num.status !== 'available' && !selectedNumbers.includes(num.id)}
            onClick={() => onSelect(num.id)}
            className={cn(
              "h-12 transition-all duration-200",
              selectedNumbers.includes(num.id) && "ring-2 ring-primary",
              num.status === 'sold' && "opacity-50 cursor-not-allowed"
            )}
          >
            {num.number}
          </Button>
        ))}
      </div>
    </div>
  );
}
```

**Resultado esperado:**
- ✅ Conexión WebSocket automática al entrar al sorteo
- ✅ Números se actualizan en tiempo real sin polling
- ✅ Auto-reconnect si se pierde la conexión
- ✅ Indicador visual de estado de conexión

---

## SEMANA 3: Sistema de Pagos Modular

### TAREA 3.1: Tabla de Payment Providers

**Migración SQL:**

```sql
-- backend/migrations/007_create_payment_providers.up.sql

CREATE TABLE payment_providers (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_enabled BOOLEAN DEFAULT false,
    is_automatic BOOLEAN DEFAULT true,
    priority INT DEFAULT 0,
    config JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Índices
CREATE INDEX idx_payment_providers_enabled ON payment_providers(is_enabled);
CREATE INDEX idx_payment_providers_priority ON payment_providers(priority);

-- Seed inicial
INSERT INTO payment_providers (code, name, description, is_enabled, is_automatic, priority, config) VALUES
(
    'bac_credomatic',
    'Tarjetas de Crédito/Débito',
    'Pago con tarjetas Visa, Mastercard vía BAC Credomatic',
    false,  -- Deshabilitado hasta tener merchant ID
    true,
    1,
    '{
        "merchant_id": "",
        "api_key": "",
        "environment": "sandbox"
    }'::jsonb
),
(
    'sinpe_movil',
    'SINPE Móvil',
    'Transferencia instantánea entre bancos costarricenses',
    true,   -- Habilitado desde el inicio
    false,  -- Requiere validación manual
    2,
    '{
        "phone_number": "8888-8888",
        "account_holder": "Sorteos CR"
    }'::jsonb
);

-- Down migration
-- backend/migrations/007_create_payment_providers.down.sql
DROP TABLE IF EXISTS payment_providers;
```

**Resultado esperado:**
- ✅ Tabla `payment_providers` creada
- ✅ 2 proveedores iniciales: BAC (deshabilitado) + SINPE (habilitado)
- ✅ Configuración en formato JSONB flexible

---

### TAREA 3.2: Domain Entity - Payment Provider

**Archivo:** `backend/internal/domain/payment_provider.go`

```go
package domain

import (
    "encoding/json"
    "errors"
    "time"
)

type PaymentProvider struct {
    ID          int64
    Code        string
    Name        string
    Description string
    IsEnabled   bool
    IsAutomatic bool
    Priority    int
    Config      json.RawMessage
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Validate valida el payment provider
func (p *PaymentProvider) Validate() error {
    if p.Code == "" {
        return errors.New("code is required")
    }
    if p.Name == "" {
        return errors.New("name is required")
    }
    if p.Priority < 0 {
        return errors.New("priority must be >= 0")
    }
    return nil
}

// Enable habilita el proveedor
func (p *PaymentProvider) Enable() {
    p.IsEnabled = true
    p.UpdatedAt = time.Now()
}

// Disable deshabilita el proveedor
func (p *PaymentProvider) Disable() {
    p.IsEnabled = false
    p.UpdatedAt = time.Now()
}

// GetConfigValue obtiene un valor de configuración
func (p *PaymentProvider) GetConfigValue(key string) (interface{}, error) {
    var config map[string]interface{}
    if err := json.Unmarshal(p.Config, &config); err != nil {
        return nil, err
    }
    
    value, ok := config[key]
    if !ok {
        return nil, errors.New("config key not found")
    }
    
    return value, nil
}

// UpdateConfig actualiza la configuración
func (p *PaymentProvider) UpdateConfig(config map[string]interface{}) error {
    configJSON, err := json.Marshal(config)
    if err != nil {
        return err
    }
    
    p.Config = configJSON
    p.UpdatedAt = time.Now()
    
    return nil
}
```

**Resultado esperado:**
- ✅ Entidad PaymentProvider con métodos de validación
- ✅ Configuración flexible en JSON
- ✅ Métodos Enable/Disable

---

### TAREA 3.3: Admin Panel - Habilitar/Deshabilitar Proveedores

**Archivo:** `backend/internal/adapters/http/handlers/admin_payment_providers.go`

```go
package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/sorteos-platform/backend/internal/domain/repositories"
)

type AdminPaymentProvidersHandler struct {
    repo repositories.PaymentProviderRepository
}

func NewAdminPaymentProvidersHandler(repo repositories.PaymentProviderRepository) *AdminPaymentProvidersHandler {
    return &AdminPaymentProvidersHandler{repo: repo}
}

// ListProviders lista todos los proveedores
func (h *AdminPaymentProvidersHandler) ListProviders(c *gin.Context) {
    providers, err := h.repo.FindAll(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching providers"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"providers": providers})
}

// UpdateProvider actualiza un proveedor
func (h *AdminPaymentProvidersHandler) UpdateProvider(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    
    var input struct {
        IsEnabled   *bool                  `json:"is_enabled"`
        Config      map[string]interface{} `json:"config"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }
    
    provider, err := h.repo.FindByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "provider not found"})
        return
    }
    
    // Actualizar campos
    if input.IsEnabled != nil {
        if *input.IsEnabled {
            provider.Enable()
        } else {
            provider.Disable()
        }
    }
    
    if input.Config != nil {
        if err := provider.UpdateConfig(input.Config); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config"})
            return
        }
    }
    
    if err := h.repo.Update(c.Request.Context(), provider); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating provider"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"provider": provider})
}
```

**Frontend Admin:**

```typescript
// frontend/src/pages/admin/PaymentProviders.tsx

import { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { Switch } from '@/components/ui/switch';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { api } from '@/lib/api';

export function PaymentProvidersPage() {
  const { data } = useQuery({
    queryKey: ['admin', 'payment-providers'],
    queryFn: () => api.get('/admin/payment-providers'),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, is_enabled }: { id: number; is_enabled: boolean }) =>
      api.patch(`/admin/payment-providers/${id}`, { is_enabled }),
  });

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Métodos de Pago</h1>
      
      {data?.providers.map((provider: any) => (
        <Card key={provider.id} className="p-4">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-semibold">{provider.name}</h3>
              <p className="text-sm text-neutral-600">{provider.description}</p>
              
              <div className="flex gap-2 mt-2">
                <Badge variant={provider.is_enabled ? 'success' : 'secondary'}>
                  {provider.is_enabled ? 'Habilitado' : 'Deshabilitado'}
                </Badge>
                <Badge variant="outline">
                  {provider.is_automatic ? 'Automático' : 'Manual'}
                </Badge>
              </div>
            </div>
            
            <Switch
              checked={provider.is_enabled}
              onCheckedChange={(checked) =>
                updateMutation.mutate({ id: provider.id, is_enabled: checked })
              }
            />
          </div>
        </Card>
      ))}
    </div>
  );
}
```

**Resultado esperado:**
- ✅ Endpoint admin para listar/actualizar proveedores
- ✅ UI para habilitar/deshabilitar con un switch
- ✅ Solo proveedores habilitados aparecen en checkout

---

### TAREA 3.4: Checkout Page con Proveedores Dinámicos

**Archivo:** `frontend/src/pages/Checkout.tsx`

```typescript
import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useParams, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Alert } from '@/components/ui/alert';
import { Clock, CreditCard } from 'lucide-react';
import { api } from '@/lib/api';
import { useTimeRemaining } from '@/hooks/useTimeRemaining';

export function CheckoutPage() {
  const { reservationId } = useParams();
  const navigate = useNavigate();
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);

  // Cargar reserva
  const { data: reservation } = useQuery({
    queryKey: ['reservation', reservationId],
    queryFn: () => api.get(`/reservations/${reservationId}`),
  });

  // Cargar proveedores habilitados
  const { data: providers } = useQuery({
    queryKey: ['payment-providers'],
    queryFn: () => api.get('/payment-providers'),
  });

  const timeRemaining = useTimeRemaining(reservation?.expires_at);
  const isUrgent = timeRemaining < 120; // < 2 minutos

  const handlePayment = async () => {
    if (!selectedProvider) return;

    // Transicionar a fase checkout si aún está en selección
    if (reservation.phase === 'selection') {
      await api.post(`/reservations/${reservationId}/move-to-checkout`);
    }

    // Redirigir según proveedor
    if (selectedProvider === 'bac_credomatic') {
      // Integración BAC
      const { payment_url } = await api.post('/payments', {
        reservation_id: reservationId,
        provider: 'bac_credomatic',
      });
      window.location.href = payment_url;
    } else if (selectedProvider === 'sinpe_movil') {
      // Flujo manual SINPE
      navigate(`/payment/sinpe/${reservationId}`);
    }
  };

  return (
    <div className="max-w-2xl mx-auto p-6 space-y-6">
      {/* Timer urgente */}
      {isUrgent && (
        <Alert variant="destructive">
          <Clock className="h-4 w-4" />
          <AlertTitle>¡Última oportunidad!</AlertTitle>
          <AlertDescription>
            Tu reserva expira en {Math.floor(timeRemaining / 60)}:
            {(timeRemaining % 60).toString().padStart(2, '0')}.
            Completa el pago o perderás los números.
          </AlertDescription>
        </Alert>
      )}

      {/* Resumen de reserva */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">Resumen de compra</h2>
        <div className="space-y-2">
          <div className="flex justify-between">
            <span>Números seleccionados:</span>
            <span className="font-semibold">{reservation?.numbers.length}</span>
          </div>
          <div className="flex justify-between">
            <span>Total:</span>
            <span className="text-2xl font-bold">
              ₡{reservation?.total_amount.toLocaleString()}
            </span>
          </div>
        </div>
      </Card>

      {/* Selección de proveedor */}
      <div className="space-y-3">
        <h3 className="font-semibold">Selecciona método de pago</h3>
        
        {providers?.filter((p: any) => p.is_enabled).map((provider: any) => (
          <Card
            key={provider.code}
            className={cn(
              "p-4 cursor-pointer transition-all",
              selectedProvider === provider.code && "ring-2 ring-primary"
            )}
            onClick={() => setSelectedProvider(provider.code)}
          >
            <div className="flex items-center gap-3">
              <CreditCard className="h-6 w-6" />
              <div>
                <p className="font-semibold">{provider.name}</p>
                <p className="text-sm text-neutral-600">{provider.description}</p>
              </div>
            </div>
          </Card>
        ))}
      </div>

      {/* Botón de pago */}
      <Button
        size="lg"
        className="w-full"
        disabled={!selectedProvider}
        onClick={handlePayment}
      >
        Pagar ₡{reservation?.total_amount.toLocaleString()}
      </Button>
    </div>
  );
}
```

**Resultado esperado:**
- ✅ Página de checkout con timer visible
- ✅ Solo proveedores habilitados se muestran
- ✅ Flujo dinámico según proveedor (BAC automático / SINPE manual)
- ✅ Transición automática a fase checkout

---

## 🧪 VALIDACIÓN Y TESTING

### Tests de Concurrencia

**Archivo:** `backend/internal/usecases/reservation_usecases_test.go`

```go
func TestReservation_ConcurrentCreation(t *testing.T) {
    // Setup
    ctx := context.Background()
    raffleID := uuid.New()
    numbers := []string{"00", "01", "02"}
    
    // 100 usuarios intentan reservar los mismos 3 números simultáneamente
    var wg sync.WaitGroup
    successCount := atomic.Int32{}
    errorCount := atomic.Int32{}
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(userIndex int) {
            defer wg.Done()
            
            _, err := useCase.CreateReservation(ctx, CreateReservationInput{
                RaffleID:  raffleID,
                UserID:    uuid.New(),
                NumberIDs: numbers,
                SessionID: fmt.Sprintf("session-%d", userIndex),
            })
            
            if err == nil {
                successCount.Add(1)
            } else {
                errorCount.Add(1)
            }
        }(i)
    }
    
    wg.Wait()
    
    // SOLO uno debe tener éxito
    assert.Equal(t, int32(1), successCount.Load())
    assert.Equal(t, int32(99), errorCount.Load())
    
    // Verificar en DB que los números están reservados
    for _, numberID := range numbers {
        number, _ := numberRepo.FindByID(ctx, numberID)
        assert.Equal(t, "reserved", number.Status)
    }
}
```

### Test de WebSocket

```bash
# Herramienta: websocat (https://github.com/vi/websocat)
websocat ws://localhost:8080/raffles/123e4567-e89b-12d3-a456-426614174000/ws

# Debería recibir mensajes cuando otros usuarios reserven números
```

### Test de Expiration Job

```go
func TestReservation_ExpireOldReservations(t *testing.T) {
    ctx := context.Background()
    
    // Crear reserva expirada (hace 15 minutos)
    reservation := &entities.Reservation{
        ExpiresAt: time.Now().Add(-15 * time.Minute),
        Status:    entities.ReservationStatusPending,
    }
    reservationRepo.Create(ctx, reservation)
    
    // Ejecutar job de expiración
    count, err := useCase.ExpireReservations(ctx)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, count)
    
    // Verificar que la reserva está expirada
    updated, _ := reservationRepo.FindByID(ctx, reservation.ID)
    assert.Equal(t, entities.ReservationStatusExpired, updated.Status)
    
    // Verificar que los números fueron liberados
    for _, numberID := range reservation.Numbers {
        number, _ := numberRepo.FindByID(ctx, numberID)
        assert.Equal(t, "available", number.Status)
    }
}
```

---

## 📊 MÉTRICAS DE ÉXITO

Al finalizar esta fase, el sistema debe cumplir:

| Métrica | Objetivo | Cómo medir |
|---------|----------|------------|
| Doble ventas | 0% | Test de concurrencia con 1000 usuarios |
| Latencia WebSocket | < 100ms | Monitor en frontend: time desde broadcast hasta render |
| Tasa de conversión | > 70% | (pagos completados / reservas creadas) * 100 |
| Uptime WebSocket | > 99% | Logs de desconexiones / reconexiones |
| P95 API latency | < 200ms | Prometheus metrics |

---

## 🚀 DEPLOYMENT CHECKLIST

Antes de desplegar a producción:

- [ ] Todos los tests unitarios pasando (coverage > 80%)
- [ ] Test de concurrencia validado con 1000 usuarios
- [ ] WebSocket probado con 100 clientes simultáneos
- [ ] Migaciones SQL aplicadas en staging
- [ ] Redis configurado con persistencia (RDB + AOF)
- [ ] Nginx configurado con WebSocket proxy (`proxy_http_version 1.1`)
- [ ] Monitoring: logs de WebSocket, métricas de locks Redis
- [ ] Variables de entorno configuradas (REDIS_URL, WS_URL, etc.)
- [ ] Job de expiración de reservas en cron (cada 1 minuto)
- [ ] Rate limiting en endpoints de reserva (10 req/min por usuario)

---

## 📚 REFERENCIAS

- **Redis Distributed Locks:** https://redis.io/docs/manual/patterns/distributed-locks/
- **Gorilla WebSocket:** https://github.com/gorilla/websocket
- **PostgreSQL Transactions:** https://www.postgresql.org/docs/current/tutorial-transactions.html
- **React WebSocket Hooks:** https://github.com/robtaussig/react-use-websocket

---

## 🆘 TROUBLESHOOTING

### Problema: Doble venta detectada

**Causa posible:** Lock no se adquirió correctamente o transacción DB falló

**Solución:**
1. Verificar logs de Redis: `redis-cli MONITOR`
2. Confirmar que TTL de locks es > duración de transacción DB
3. Revisar logs de PostgreSQL para deadlocks

### Problema: WebSocket desconecta frecuentemente

**Causa posible:** Nginx timeout o red inestable

**Solución:**
1. Aumentar timeout en Nginx:
```nginx
proxy_read_timeout 3600s;
proxy_send_timeout 3600s;
```
2. Implementar ping/pong (ya incluido en Client.WritePump)

### Problema: Números no se liberan al expirar

**Causa posible:** Job de expiración no está corriendo

**Solución:**
1. Verificar que cron job está activo
2. Revisar logs del job
3. Ejecutar manualmente: `curl -X POST http://localhost:8080/admin/jobs/expire-reservations`

---

**Fin del documento**

Última actualización: 2025-11-12  
Responsable: Equipo de Desarrollo  
Revisado por: Ing. Alonso Alpízar