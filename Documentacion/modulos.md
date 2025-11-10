# Módulos del Sistema - Plataforma de Sorteos

**Versión:** 1.0
**Fecha:** 2025-11-10
**Arquitectura:** Hexagonal / Clean Architecture

---

## 1. Introducción

Este documento describe los **módulos obligatorios** del sistema, su responsabilidad, interfaces, dependencias y flujos críticos. La arquitectura sigue el patrón hexagonal para garantizar:

- **Independencia del dominio**: Reglas de negocio sin dependencias externas
- **Testabilidad**: Cada capa es testeable de forma aislada
- **Extensibilidad**: Nuevos adapters (PSPs, notificadores) sin cambiar el core

---

## 2. Arquitectura de Capas

```
┌─────────────────────────────────────────────────────┐
│               HTTP Handlers (Gin)                    │
│          (Adapters/Driving - Entradas)              │
└─────────────────┬───────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────┐
│                Use Cases                             │
│         (Application Layer - Lógica de App)         │
│  CreateRaffle, ReserveNumbers, ProcessPayment, etc. │
└─────────────────┬───────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────┐
│                  Domain                              │
│         (Entities + Business Rules)                  │
│     Raffle, User, Reservation, Payment, etc.        │
└─────────────────┬───────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────┐
│             Adapters (Driven - Salidas)             │
│  Repositories (GORM), PaymentProviders, Notifiers   │
└─────────────────────────────────────────────────────┘
```

---

## 3. Módulo 1: Auth & Perfil

### 3.1 Responsabilidad

- Registro de usuarios (email/teléfono)
- Verificación de identidad (email, SMS, KYC)
- Autenticación (login, refresh token)
- Gestión de perfil (datos básicos, direcciones, medios de pago)
- Control de acceso basado en roles (RBAC)

---

### 3.2 Entidades de Dominio

#### User
```go
type User struct {
    ID           int64
    Email        string // unique, validado
    Phone        string // +506XXXXXXXX (formato E.164)
    PasswordHash string
    Role         UserRole // user, admin
    KYCLevel     KYCLevel // none, email_verified, phone_verified, full_kyc
    Status       UserStatus // active, suspended, banned
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type UserRole string
const (
    RoleUser  UserRole = "user"
    RoleAdmin UserRole = "admin"
)

type KYCLevel string
const (
    KYCNone          KYCLevel = "none"
    KYCEmailVerified KYCLevel = "email_verified"
    KYCPhoneVerified KYCLevel = "phone_verified"
    KYCFullVerified  KYCLevel = "full_kyc" // requiere ID, selfie, etc.
)
```

#### PaymentMethod
```go
type PaymentMethod struct {
    ID         int64
    UserID     int64
    Type       PaymentMethodType // card, paypal, bank_account
    Provider   string // stripe, paypal
    ExternalID string // tok_xxx (Stripe), ba_xxx (PayPal)
    Last4      string
    IsDefault  bool
    CreatedAt  time.Time
}
```

---

### 3.3 Casos de Uso (Use Cases)

#### RegisterUser
**Input:**
```go
type RegisterUserInput struct {
    Email    string
    Phone    string
    Password string
}
```

**Flujo:**
1. Validar formato de email/phone
2. Hash de contraseña con bcrypt
3. Crear usuario en DB (status=active, kyc_level=none)
4. Generar token de verificación (JWT de corta vida)
5. Enviar email/SMS de verificación
6. Retornar user_id

**Salida:**
```go
type RegisterUserOutput struct {
    UserID int64
    Message string // "Revisa tu email para verificar tu cuenta"
}
```

---

#### VerifyEmail
**Input:**
```go
type VerifyEmailInput struct {
    Token string // JWT con claim user_id + exp
}
```

**Flujo:**
1. Decodificar y validar JWT
2. Actualizar user.kyc_level = email_verified
3. Generar access_token + refresh_token
4. Retornar tokens

---

#### Login
**Input:**
```go
type LoginInput struct {
    Email    string
    Password string
}
```

**Flujo:**
1. Buscar usuario por email
2. Comparar password hash
3. Validar user.status == active
4. Generar access_token (15 min) + refresh_token (7 días)
5. Guardar refresh_token en Redis con TTL
6. Retornar tokens + user info

**Salida:**
```go
type LoginOutput struct {
    AccessToken  string
    RefreshToken string
    User         UserDTO
}
```

---

#### RefreshToken
**Input:**
```go
type RefreshTokenInput struct {
    RefreshToken string
}
```

**Flujo:**
1. Validar token en Redis
2. Decodificar y verificar claims
3. Generar nuevo access_token
4. Rotar refresh_token (invalidar anterior, crear nuevo)
5. Retornar nuevos tokens

---

### 3.4 Endpoints HTTP

| Método | Ruta | Handler | Auth |
|--------|------|---------|------|
| POST | /auth/register | RegisterHandler | No |
| POST | /auth/verify-email | VerifyEmailHandler | No |
| POST | /auth/login | LoginHandler | No |
| POST | /auth/refresh | RefreshTokenHandler | No |
| POST | /auth/logout | LogoutHandler | Sí |
| GET | /users/me | GetProfileHandler | Sí |
| PATCH | /users/me | UpdateProfileHandler | Sí |
| POST | /users/me/payment-methods | AddPaymentMethodHandler | Sí |

---

### 3.5 Interfaces (Ports)

```go
// Repository
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByID(ctx context.Context, id int64) (*User, error)
    Update(ctx context.Context, user *User) error
}

// Token Manager
type TokenManager interface {
    GenerateAccessToken(userID int64, role UserRole) (string, error)
    GenerateRefreshToken(userID int64) (string, error)
    ValidateToken(token string) (*TokenClaims, error)
    RevokeRefreshToken(ctx context.Context, token string) error
}

// Notifier
type Notifier interface {
    SendVerificationEmail(ctx context.Context, to string, token string) error
    SendVerificationSMS(ctx context.Context, to string, code string) error
}
```

---

## 4. Módulo 2: Publicación de Sorteos

### 4.1 Responsabilidad

- Crear/editar/suspender sorteos
- Gestionar detalles (título, descripción, imágenes, rango de números)
- Configurar fuente de sorteo (Lotería CR, fecha específica)
- Listar sorteos con filtros y paginación
- Validar reglas de publicación (parámetros dinámicos)

---

### 4.2 Entidades de Dominio

#### Raffle
```go
type Raffle struct {
    ID             int64
    UserID         int64 // owner
    Title          string
    Description    string
    Status         RaffleStatus
    DrawDate       time.Time
    LotterySource  string // "loteria_nacional_cr"
    LotteryDate    string // "2025-12-25"
    NumberRange    NumberRange // ej: {Min: 0, Max: 99}
    PricePerNumber decimal.Decimal
    TotalNumbers   int
    SoldCount      int
    ReservedCount  int
    CreatedAt      time.Time
    PublishedAt    *time.Time
}

type RaffleStatus string
const (
    RaffleStatusDraft     RaffleStatus = "draft"
    RaffleStatusActive    RaffleStatus = "active"
    RaffleStatusSuspended RaffleStatus = "suspended"
    RaffleStatusCompleted RaffleStatus = "completed"
    RaffleStatusCancelled RaffleStatus = "cancelled"
)

type NumberRange struct {
    Min int // 0
    Max int // 99 (genera 00, 01, ..., 99)
}
```

#### RaffleNumber
```go
type RaffleNumber struct {
    RaffleID   int64
    Number     string // "00" a "99" (con ceros a la izquierda)
    UserID     *int64 // null si no vendido
    Status     NumberStatus
    ReservedAt *time.Time
    SoldAt     *time.Time
}

type NumberStatus string
const (
    NumberAvailable NumberStatus = "available"
    NumberReserved  NumberStatus = "reserved"
    NumberSold      NumberStatus = "sold"
)
```

#### RaffleImage
```go
type RaffleImage struct {
    ID       int64
    RaffleID int64
    URL      string // https://cdn.example.com/raffles/123/img1.jpg
    Order    int
}
```

---

### 4.3 Casos de Uso

#### CreateRaffle
**Input:**
```go
type CreateRaffleInput struct {
    UserID         int64
    Title          string
    Description    string
    DrawDate       time.Time
    LotterySource  string
    LotteryDate    string
    NumberRangeMin int
    NumberRangeMax int
    PricePerNumber decimal.Decimal
    Images         []ImageUpload
}
```

**Flujo:**
1. Validar parámetros (ver parametrizacion_reglas.md)
   - DrawDate debe ser futuro
   - PricePerNumber > 0
   - NumberRange válido (Max - Min + 1 = TotalNumbers)
2. Crear raffle en DB (status=draft)
3. Generar números (ej: 00-99) → insertar en raffle_numbers
4. Subir imágenes a S3/storage → insertar en raffle_images
5. Retornar raffle_id

**Reglas de Negocio:**
- Usuario debe tener KYC >= email_verified
- Máximo 10 sorteos activos por usuario (parámetro configurable)

---

#### PublishRaffle
**Input:**
```go
type PublishRaffleInput struct {
    RaffleID int64
    UserID   int64 // validar ownership
}
```

**Flujo:**
1. Verificar raffle.status == draft
2. Validar reglas pre-publicación (ver parametrizacion_reglas.md)
3. Actualizar status = active, published_at = now()
4. Enviar notificaciones a seguidores (futuro)
5. Invalidar caché Redis de listados

---

#### SuspendRaffle (Admin)
**Input:**
```go
type SuspendRaffleInput struct {
    RaffleID int64
    AdminID  int64
    Reason   string
}
```

**Flujo:**
1. Validar admin.role == admin
2. Actualizar raffle.status = suspended
3. Registrar en audit_logs (admin_id, raffle_id, reason)
4. Notificar al owner
5. Rechazar nuevas reservas en este sorteo

---

#### ListRaffles
**Input:**
```go
type ListRafflesInput struct {
    Status    *RaffleStatus
    UserID    *int64 // filtrar por owner
    Category  *string
    MinPrice  *decimal.Decimal
    MaxPrice  *decimal.Decimal
    Page      int
    PageSize  int
    SortBy    string // "created_at", "draw_date", "price"
    SortOrder string // "asc", "desc"
}
```

**Flujo:**
1. Construir query dinámica con filtros
2. Aplicar paginación (OFFSET/LIMIT)
3. Cachear resultado en Redis (key: hash de filtros, TTL: 5 min)
4. Retornar lista + metadata (total_count, total_pages)

---

### 4.4 Endpoints HTTP

| Método | Ruta | Handler | Auth | Rol |
|--------|------|---------|------|-----|
| POST | /raffles | CreateRaffleHandler | Sí | user |
| GET | /raffles | ListRafflesHandler | No | - |
| GET | /raffles/:id | GetRaffleDetailHandler | No | - |
| PATCH | /raffles/:id | UpdateRaffleHandler | Sí | owner |
| POST | /raffles/:id/publish | PublishRaffleHandler | Sí | owner |
| PATCH | /admin/raffles/:id/suspend | SuspendRaffleHandler | Sí | admin |

---

### 4.5 Interfaces

```go
type RaffleRepository interface {
    Create(ctx context.Context, raffle *Raffle) error
    FindByID(ctx context.Context, id int64) (*Raffle, error)
    Update(ctx context.Context, raffle *Raffle) error
    List(ctx context.Context, filters RaffleFilters) ([]*Raffle, int, error)
    GenerateNumbers(ctx context.Context, raffleID int64, numberRange NumberRange) error
}

type ImageStorage interface {
    Upload(ctx context.Context, file io.Reader, filename string) (string, error) // retorna URL
    Delete(ctx context.Context, url string) error
}
```

---

## 5. Módulo 3: Reserva y Compra de Números

### 5.1 Responsabilidad

- Reservar números temporalmente (5 min)
- Crear pago con PSP
- Confirmar venta tras pago exitoso
- Liberar reservas expiradas o pagos fallidos
- Evitar doble venta con locks distribuidos

---

### 5.2 Entidades de Dominio

#### Reservation
```go
type Reservation struct {
    ID             int64
    RaffleID       int64
    UserID         int64
    Numbers        []string // ["01", "15", "42"]
    Status         ReservationStatus
    IdempotencyKey string // UUID generado por cliente
    ExpiresAt      time.Time // now() + 5 min
    CreatedAt      time.Time
}

type ReservationStatus string
const (
    ReservationPending   ReservationStatus = "pending"
    ReservationConfirmed ReservationStatus = "confirmed"
    ReservationExpired   ReservationStatus = "expired"
    ReservationCancelled ReservationStatus = "cancelled"
)
```

---

### 5.3 Casos de Uso

#### ReserveNumbers (Crítico - Alta Concurrencia)

**Input:**
```go
type ReserveNumbersInput struct {
    RaffleID       int64
    UserID         int64
    Numbers        []string
    IdempotencyKey string
}
```

**Flujo (con locks distribuidos):**

```go
func (uc *ReserveNumbersUseCase) Execute(ctx context.Context, input ReserveNumbersInput) (*Reservation, error) {
    // 1. Validaciones previas
    raffle, err := uc.raffleRepo.FindByID(ctx, input.RaffleID)
    if raffle.Status != RaffleStatusActive {
        return nil, ErrRaffleNotActive
    }

    // 2. Verificar idempotencia
    if existing := uc.reservationRepo.FindByIdempotencyKey(ctx, input.IdempotencyKey); existing != nil {
        return existing, nil
    }

    // 3. Adquirir locks distribuidos (Redis)
    locks := []string{}
    for _, num := range input.Numbers {
        lockKey := fmt.Sprintf("lock:raffle:%d:num:%s", input.RaffleID, num)
        acquired, err := uc.redis.SetNX(ctx, lockKey, input.UserID, 30*time.Second).Result()
        if !acquired {
            // Liberar locks ya adquiridos
            for _, l := range locks {
                uc.redis.Del(ctx, l)
            }
            return nil, ErrNumberAlreadyReserved
        }
        locks = append(locks, lockKey)
    }
    defer func() {
        for _, l := range locks {
            uc.redis.Del(ctx, l)
        }
    }()

    // 4. Crear reserva en transacción
    reservation := &Reservation{
        RaffleID:       input.RaffleID,
        UserID:         input.UserID,
        Numbers:        input.Numbers,
        Status:         ReservationPending,
        IdempotencyKey: input.IdempotencyKey,
        ExpiresAt:      time.Now().Add(5 * time.Minute),
    }

    err = uc.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(reservation).Error; err != nil {
            return err
        }
        // Actualizar raffle_numbers a reserved
        return tx.Model(&RaffleNumber{}).
            Where("raffle_id = ? AND number IN ? AND status = ?",
                input.RaffleID, input.Numbers, NumberAvailable).
            Updates(map[string]interface{}{
                "status": NumberReserved,
                "user_id": input.UserID,
                "reserved_at": time.Now(),
            }).Error
    })

    if err != nil {
        return nil, err
    }

    // 5. Guardar en Redis para TTL automático
    uc.redis.Set(ctx,
        fmt.Sprintf("reservation:%d", reservation.ID),
        reservation,
        5*time.Minute,
    )

    return reservation, nil
}
```

**Reglas de Negocio:**
- Máximo 10 números por reserva (parámetro configurable)
- Usuario debe tener KYC >= email_verified
- Números deben estar available (no reserved ni sold)

---

#### ReleaseExpiredReservations (Cron Job)

**Flujo:**
```go
func (uc *ReleaseExpiredReservationsUseCase) Execute(ctx context.Context) error {
    expiredReservations := uc.reservationRepo.FindExpired(ctx)

    for _, res := range expiredReservations {
        uc.db.Transaction(func(tx *gorm.DB) error {
            // Actualizar reserva a expired
            tx.Model(&Reservation{}).Where("id = ?", res.ID).
                Update("status", ReservationExpired)

            // Liberar números
            tx.Model(&RaffleNumber{}).
                Where("raffle_id = ? AND number IN ?", res.RaffleID, res.Numbers).
                Updates(map[string]interface{}{
                    "status": NumberAvailable,
                    "user_id": nil,
                    "reserved_at": nil,
                })

            return nil
        })
    }

    return nil
}
```

**Ejecución:** Cada 1 minuto (cron)

---

### 5.4 Endpoints HTTP

| Método | Ruta | Handler | Auth |
|--------|------|---------|------|
| POST | /raffles/:id/reservations | CreateReservationHandler | Sí |
| GET | /reservations/:id | GetReservationHandler | Sí |
| DELETE | /reservations/:id | CancelReservationHandler | Sí |

---

### 5.5 Interfaces

```go
type ReservationRepository interface {
    Create(ctx context.Context, reservation *Reservation) error
    FindByID(ctx context.Context, id int64) (*Reservation, error)
    FindByIdempotencyKey(ctx context.Context, key string) (*Reservation, error)
    FindExpired(ctx context.Context) ([]*Reservation, error)
    Update(ctx context.Context, reservation *Reservation) error
}

type LockManager interface {
    AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
    ReleaseLock(ctx context.Context, key string) error
}
```

---

## 6. Módulo 4: Procesamiento de Pagos

**Ver:** [pagos_integraciones.md](./pagos_integraciones.md) para detalles completos.

### 6.1 Responsabilidad

- Integrar con múltiples PSPs (Stripe, PayPal, local)
- Procesar pagos con idempotencia
- Manejar webhooks de confirmación/fallo
- Gestionar chargebacks y reembolsos

---

### 6.2 Entidades

#### Payment
```go
type Payment struct {
    ID             int64
    ReservationID  int64
    UserID         int64
    Provider       string // "stripe", "paypal"
    Amount         decimal.Decimal
    Currency       string // "USD", "CRC"
    Status         PaymentStatus
    ExternalID     string // pi_xxx (Stripe), PAYID-xxx (PayPal)
    IdempotencyKey string
    Metadata       map[string]interface{} // JSONB
    CreatedAt      time.Time
}

type PaymentStatus string
const (
    PaymentPending   PaymentStatus = "pending"
    PaymentSucceeded PaymentStatus = "succeeded"
    PaymentFailed    PaymentStatus = "failed"
    PaymentRefunded  PaymentStatus = "refunded"
)
```

---

### 6.3 Interfaz PaymentProvider

```go
type PaymentProvider interface {
    Authorize(ctx context.Context, input AuthorizeInput) (*AuthorizeOutput, error)
    Capture(ctx context.Context, paymentID string) error
    Refund(ctx context.Context, paymentID string, amount decimal.Decimal) error
    VerifyWebhook(ctx context.Context, payload []byte, signature string) (*WebhookEvent, error)
}
```

---

## 7. Módulo 5: Selección de Ganador

### 7.1 Responsabilidad

- Consultar API de lotería en draw_date
- Determinar ganadores según fuente oficial
- Notificar a ganadores y owners
- Iniciar proceso de liquidación

---

### 7.2 Caso de Uso: DetermineWinner

**Flujo:**
```go
func (uc *DetermineWinnerUseCase) Execute(ctx context.Context, raffleID int64) error {
    raffle := uc.raffleRepo.FindByID(ctx, raffleID)

    // 1. Consultar lotería
    result, err := uc.lotteryClient.GetResult(raffle.LotterySource, raffle.LotteryDate)
    if err != nil {
        return err
    }

    // 2. Extraer número ganador (ej: últimos 2 dígitos)
    winningNumber := extractNumber(result.WinningNumber)

    // 3. Buscar ganador
    winner, err := uc.raffleNumberRepo.FindByRaffleAndNumber(ctx, raffleID, winningNumber)
    if winner.UserID == nil {
        // Número no vendido, no hay ganador
        raffle.WinnerID = nil
    } else {
        raffle.WinnerID = winner.UserID
    }

    raffle.Status = RaffleStatusCompleted
    raffle.WinningNumber = winningNumber
    uc.raffleRepo.Update(ctx, raffle)

    // 4. Notificar
    if raffle.WinnerID != nil {
        uc.notifier.SendEmail(ctx, winner.Email, "winner_notification", map[string]any{
            "raffle_title": raffle.Title,
            "number": winningNumber,
        })
    }

    return nil
}
```

---

## 8. Módulo 6: Backoffice (Admin)

### 8.1 Responsabilidad

- Gestionar sorteos (suspender/activar/eliminar)
- Gestionar usuarios (verificar KYC, suspender)
- Crear liquidaciones manuales
- Ver logs de auditoría
- Generar reportes

---

### 8.2 Entidades

#### AuditLog
```go
type AuditLog struct {
    ID         int64
    UserID     int64 // admin que realizó la acción
    Action     string // "suspend_raffle", "verify_kyc"
    EntityType string // "raffle", "user"
    EntityID   int64
    IPAddress  string
    UserAgent  string
    Metadata   map[string]interface{} // JSONB con detalles
    CreatedAt  time.Time
}
```

---

### 8.3 Endpoints

| Método | Ruta | Handler | Rol |
|--------|------|---------|-----|
| GET | /admin/raffles | ListAllRafflesHandler | admin |
| PATCH | /admin/raffles/:id | UpdateRaffleAdminHandler | admin |
| GET | /admin/users | ListUsersHandler | admin |
| PATCH | /admin/users/:id/kyc | VerifyKYCHandler | admin |
| POST | /admin/settlements | CreateSettlementHandler | admin |
| GET | /admin/audit-logs | ListAuditLogsHandler | admin |

---

## 9. Módulo 7: Notificaciones

### 9.1 Responsabilidad

- Enviar emails transaccionales (verificación, pago, ganador)
- Enviar SMS (verificación de teléfono)
- Notificaciones push (futuro - Fase 3)
- Webhooks a sistemas externos (futuro)

---

### 9.2 Interfaz Notifier

```go
type Notifier interface {
    SendEmail(ctx context.Context, to string, template string, data map[string]interface{}) error
    SendSMS(ctx context.Context, to string, message string) error
}
```

**Templates de email:**
- `verification_email`
- `reservation_confirmed`
- `payment_succeeded`
- `raffle_winner`
- `raffle_completed_owner`

---

## 10. Dependencias entre Módulos

```
Auth & Perfil
    ↓ (UserID)
Publicación de Sorteos
    ↓ (RaffleID)
Reserva y Compra
    ↓ (ReservationID)
Procesamiento de Pagos
    ↓ (PaymentID)
Selección de Ganador
    ↓ (WinnerID)
Liquidaciones (Backoffice)
```

---

## 11. Tests Críticos por Módulo

### Auth
- [ ] Registro con email duplicado falla
- [ ] Token expirado es rechazado
- [ ] Refresh token rotation funciona correctamente

### Reservas
- [ ] 500 solicitudes concurrentes no generan doble venta
- [ ] Reserva expira a los 5 min exactos
- [ ] Lock distribuido se libera si handler crashea

### Pagos
- [ ] Mismo Idempotency-Key retorna pago existente
- [ ] Webhook tardío no confirma reserva expirada
- [ ] Refund marca números como disponibles

### Ganadores
- [ ] Sorteo sin número ganador vendido marca winner_id = null
- [ ] Notificación se envía solo si hay ganador

---

## 12. Próximos Pasos

1. Diseñar esquema de base de datos completo
2. Implementar repositorios con GORM
3. Escribir tests unitarios para use cases
4. Implementar handlers HTTP con Gin
5. Integrar con Redis para locks y caché

---

**Ver también:**
- [stack_tecnico.md](./stack_tecnico.md) - Tecnologías y librerías
- [pagos_integraciones.md](./pagos_integraciones.md) - Detalles de pagos
- [seguridad.md](./seguridad.md) - Consideraciones de seguridad
