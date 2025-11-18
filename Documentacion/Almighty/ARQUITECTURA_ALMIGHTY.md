# Arquitectura del Módulo Almighty Admin

**Versión:** 1.0
**Fecha:** 2025-11-18
**Sistema:** Sorteos.club - Plataforma de Rifas

---

## 1. Visión General

El módulo **Almighty Admin** se integra en el sistema existente de Sorteos.club siguiendo los mismos principios de **Arquitectura Hexagonal (Puertos y Adaptadores)** y **Clean Architecture**.

### 1.1 Principios Arquitectónicos

- **Separación de Responsabilidades:** Lógica de negocio independiente de frameworks
- **Inversión de Dependencias:** Dominio no depende de infraestructura
- **Testabilidad:** Cada capa es testeable independientemente
- **Escalabilidad:** Diseño preparado para crecimiento
- **Seguridad:** Protección en todas las capas

---

## 2. Diagrama de Capas

```
┌─────────────────────────────────────────────────────────────────┐
│                         PRESENTATION LAYER                       │
│                   (Frontend - React + TypeScript)                │
├─────────────────────────────────────────────────────────────────┤
│  /admin/*  ┌────────────┐  ┌────────────┐  ┌────────────┐      │
│            │  Dashboard │  │   Users    │  │ Organizers │      │
│            │    Page    │  │    Page    │  │    Page    │      │
│            └────────────┘  └────────────┘  └────────────┘      │
│            ┌────────────┐  ┌────────────┐  ┌────────────┐      │
│            │  Raffles   │  │ Settlements│  │  Reports   │      │
│            │    Page    │  │    Page    │  │    Page    │      │
│            └────────────┘  └────────────┘  └────────────┘      │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  │ HTTPS (Nginx SSL)
                                  ↓
┌─────────────────────────────────────────────────────────────────┐
│                        APPLICATION LAYER                         │
│                 (HTTP Handlers - Gin Framework)                  │
├─────────────────────────────────────────────────────────────────┤
│  /api/v1/admin/*                                                 │
│                                                                   │
│  ┌────────────────────────────────────────────────────────┐     │
│  │  Middleware Chain                                       │     │
│  │  1. Authenticate() - Validar JWT                        │     │
│  │  2. RequireRole("super_admin") - Verificar permisos     │     │
│  │  3. RateLimiter (10 req/min) - Prevenir abuso          │     │
│  │  4. AuditLogger - Registrar todas las acciones         │     │
│  └────────────────────────────────────────────────────────┘     │
│                                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │UserHandler   │  │OrganizerHdlr │  │RaffleHandler │          │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤          │
│  │List()        │  │List()        │  │List()        │          │
│  │GetByID()     │  │GetByID()     │  │ForceStatus() │          │
│  │UpdateStatus()│  │SetCommission│  │CancelRefund()│          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ↓
┌─────────────────────────────────────────────────────────────────┐
│                         USE CASE LAYER                           │
│                     (Business Logic - Go)                        │
├─────────────────────────────────────────────────────────────────┤
│  internal/usecase/admin/                                         │
│                                                                   │
│  user/                     organizer/              raffle/       │
│  ├─ ListUsersUseCase       ├─ ListOrganizers      ├─ ListAdmin  │
│  ├─ UpdateStatusUseCase    ├─ SetCommission       ├─ ForceStatus│
│  ├─ UpdateKYCUseCase       └─ CalculateRevenue    └─ ManualDraw │
│  └─ ResetPasswordUseCase                                         │
│                                                                   │
│  payment/                  settlement/             reports/      │
│  ├─ ListPaymentsAdmin      ├─ CreateSettlement    ├─ Dashboard  │
│  ├─ ProcessRefundUseCase   ├─ ApproveSettlement   ├─ Revenue    │
│  └─ ManageDisputeUseCase   └─ MarkPaidUseCase     └─ Export     │
│                                                                   │
│  system/                   category/                             │
│  ├─ UpdateParameterUseCase ├─ CreateCategoryUC                  │
│  ├─ GetCompanySettings     ├─ UpdateCategoryUC                  │
│  └─ UpdatePaymentProcessor └─ ReorderCategoriesUC               │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ↓
┌─────────────────────────────────────────────────────────────────┐
│                          DOMAIN LAYER                            │
│                  (Entities & Interfaces - Go)                    │
├─────────────────────────────────────────────────────────────────┤
│  internal/domain/                                                │
│                                                                   │
│  Entities (Agregados de Negocio):                               │
│  ┌──────────────────┐  ┌──────────────────┐                    │
│  │ User             │  │ OrganizerProfile │                    │
│  ├──────────────────┤  ├──────────────────┤                    │
│  │ - ID             │  │ - ID             │                    │
│  │ - Email          │  │ - UserID         │                    │
│  │ - Role           │  │ - BusinessName   │                    │
│  │ - KYCLevel       │  │ - BankInfo       │                    │
│  │ - Status         │  │ - Commission%    │                    │
│  │ - SuspendedBy    │  │ - TotalPayouts   │                    │
│  │ - SuspendedAt    │  │ - PendingPayout  │                    │
│  └──────────────────┘  └──────────────────┘                    │
│                                                                   │
│  ┌──────────────────┐  ┌──────────────────┐                    │
│  │ Settlement       │  │ SystemParameter  │                    │
│  ├──────────────────┤  ├──────────────────┤                    │
│  │ - ID             │  │ - Key            │                    │
│  │ - RaffleID       │  │ - Value          │                    │
│  │ - OrganizerID    │  │ - ValueType      │                    │
│  │ - GrossRevenue   │  │ - Category       │                    │
│  │ - PlatformFee    │  │ - IsSensitive    │                    │
│  │ - NetPayout      │  │ - UpdatedBy      │                    │
│  │ - Status         │  └──────────────────┘                    │
│  │ - ApprovedBy     │                                            │
│  └──────────────────┘  ┌──────────────────┐                    │
│                         │ CompanySettings  │                    │
│  ┌──────────────────┐  ├──────────────────┤                    │
│  │ PaymentProcessor │  │ - CompanyName    │                    │
│  ├──────────────────┤  │ - TaxID          │                    │
│  │ - ID             │  │ - Address        │                    │
│  │ - Provider       │  │ - ContactInfo    │                    │
│  │ - Name           │  │ - LogoURL        │                    │
│  │ - IsActive       │  └──────────────────┘                    │
│  │ - Credentials    │                                            │
│  └──────────────────┘                                            │
│                                                                   │
│  Repository Interfaces (Puertos):                               │
│  - UserRepository                                                │
│  - OrganizerProfileRepository                                   │
│  - SettlementRepository                                          │
│  - SystemParameterRepository                                     │
│  - PaymentProcessorRepository                                    │
│  - CompanySettingsRepository                                     │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ↓
┌─────────────────────────────────────────────────────────────────┐
│                      INFRASTRUCTURE LAYER                        │
│                     (Adaptadores - Go + SQL)                     │
├─────────────────────────────────────────────────────────────────┤
│  internal/adapters/                                              │
│                                                                   │
│  db/ (PostgreSQL Repositories)                                   │
│  ┌─────────────────────────────────────────────────────┐        │
│  │ PostgresUserRepository                              │        │
│  │ PostgresOrganizerProfileRepository                  │        │
│  │ PostgresSettlementRepository                        │        │
│  │ PostgresSystemParameterRepository                   │        │
│  │ PostgresPaymentProcessorRepository                  │        │
│  │ PostgresCompanySettingsRepository                   │        │
│  └─────────────────────────────────────────────────────┘        │
│                                                                   │
│  redis/ (Cache & Distributed Locks)                              │
│  ┌─────────────────────────────────────────────────────┐        │
│  │ TokenManager (JWT blacklist)                        │        │
│  │ RateLimiter (admin endpoints)                       │        │
│  │ CacheService (dashboard KPIs, reports)              │        │
│  └─────────────────────────────────────────────────────┘        │
│                                                                   │
│  payments/ (Payment Providers)                                   │
│  ┌─────────────────────────────────────────────────────┐        │
│  │ StripeProvider (refunds, disputes)                  │        │
│  │ PayPalProvider (refunds, disputes)                  │        │
│  └─────────────────────────────────────────────────────┘        │
│                                                                   │
│  notifier/ (Email & SMS)                                         │
│  ┌─────────────────────────────────────────────────────┐        │
│  │ EmailService (notificaciones a usuarios/admins)     │        │
│  │ Templates: user_suspended, settlement_approved, etc. │        │
│  └─────────────────────────────────────────────────────┘        │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ↓
┌─────────────────────────────────────────────────────────────────┐
│                        PERSISTENCE LAYER                         │
├─────────────────────────────────────────────────────────────────┤
│  PostgreSQL 16                   Redis 7                         │
│  ┌────────────────────┐          ┌────────────────────┐         │
│  │ Existing Tables:   │          │ Keys:              │         │
│  │ - users            │          │ - token:blacklist  │         │
│  │ - raffles          │          │ - ratelimit:admin  │         │
│  │ - payments         │          │ - cache:dashboard  │         │
│  │ - audit_logs       │          └────────────────────┘         │
│  │                    │                                          │
│  │ New Tables:        │                                          │
│  │ - company_settings │                                          │
│  │ - payment_procs    │                                          │
│  │ - organizer_profs  │                                          │
│  │ - settlements      │                                          │
│  │ - system_params    │                                          │
│  └────────────────────┘                                          │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. Flujo de Datos: Ejemplo de Suspender Usuario

```
1. FRONTEND
   ┌─────────────────────────────────────────────┐
   │ UserDetailPage.tsx                          │
   │ - Admin hace clic en "Suspender Usuario"   │
   │ - Modal pide razón de suspensión            │
   │ - Llama a hook: useUpdateUserStatus()       │
   └─────────────────────────────────────────────┘
                      │
                      │ POST /api/v1/admin/users/:id/status
                      │ { status: "suspended", reason: "..." }
                      │ Authorization: Bearer <JWT>
                      ↓
2. HTTP HANDLER
   ┌─────────────────────────────────────────────┐
   │ UserHandler.UpdateStatus(c *gin.Context)    │
   │ 1. Middleware valida JWT                    │
   │ 2. Middleware verifica role = super_admin   │
   │ 3. Rate limiter verifica límite             │
   │ 4. Parse request body                       │
   │ 5. Validar inputs (user_id, status, reason) │
   └─────────────────────────────────────────────┘
                      │
                      │ Llama a use case
                      ↓
3. USE CASE
   ┌─────────────────────────────────────────────┐
   │ UpdateUserStatusUseCase.Execute()           │
   │ 1. Validar que status sea válido            │
   │ 2. Validar que admin no se suspenda a sí    │
   │ 3. Obtener user actual del repo             │
   │ 4. Actualizar campos:                       │
   │    - status = "suspended"                   │
   │    - suspended_by = admin_id                │
   │    - suspended_at = now()                   │
   │    - suspension_reason = "..."              │
   │ 5. Guardar en repo                          │
   │ 6. Crear audit log                          │
   │ 7. Enviar email al usuario                  │
   └─────────────────────────────────────────────┘
                      │
                      │ user_repo.Update(user)
                      ↓
4. REPOSITORY
   ┌─────────────────────────────────────────────┐
   │ PostgresUserRepository.Update(user)         │
   │ UPDATE users SET                            │
   │   status = $1,                              │
   │   suspended_by = $2,                        │
   │   suspended_at = $3,                        │
   │   suspension_reason = $4                    │
   │ WHERE id = $5                               │
   └─────────────────────────────────────────────┘
                      │
                      │ SQL Transaction
                      ↓
5. DATABASE
   ┌─────────────────────────────────────────────┐
   │ PostgreSQL - tabla users                    │
   │ Row actualizada con status = "suspended"    │
   └─────────────────────────────────────────────┘
                      │
                      │ Commit exitoso
                      ↓
6. AUDIT LOG
   ┌─────────────────────────────────────────────┐
   │ audit_logs table                            │
   │ INSERT INTO audit_logs (                    │
   │   admin_id,                                 │
   │   action = "user_suspended",                │
   │   severity = "warning",                     │
   │   entity_type = "user",                     │
   │   entity_id = user_id,                      │
   │   metadata = {reason, status_before}        │
   │ )                                           │
   └─────────────────────────────────────────────┘
                      │
                      │ Return success
                      ↓
7. FRONTEND
   ┌─────────────────────────────────────────────┐
   │ Toast: "Usuario suspendido exitosamente"    │
   │ Refresh user detail                         │
   │ Badge actualizado a "Suspended"             │
   └─────────────────────────────────────────────┘
```

---

## 4. Integración con Sistema Existente

### 4.1 Tablas Existentes Utilizadas

El módulo Almighty **NO crea tablas duplicadas**, sino que extiende y utiliza las existentes:

| Tabla Existente | Uso en Almighty | Modificaciones |
|-----------------|-----------------|----------------|
| `users` | Gestión de usuarios, filtrado, suspensión | **ALTER:** agregar suspended_by, suspension_reason, last_kyc_review, kyc_reviewer |
| `raffles` | Gestión admin de rifas, suspensión, cancelación | **ALTER:** agregar suspended_by, suspension_reason, suspended_at, admin_notes |
| `payments` | Procesar refunds, ver transacciones | **Ninguna** (usa campos existentes) |
| `reservations` | Ver reservas al cancelar rifas | **Ninguna** |
| `audit_logs` | Registrar todas las acciones de admin | **Ninguna** (usa actions existentes + nuevos) |
| `categories` | CRUD de categorías, reordenamiento | **Ninguna** |

### 4.2 Tablas Nuevas Creadas

| Tabla Nueva | Propósito |
|-------------|-----------|
| `company_settings` | Datos maestros de la empresa Sorteos.club |
| `payment_processors` | Configuración de Stripe, PayPal, etc. |
| `organizer_profiles` | Perfiles extendidos de organizadores con info bancaria y comisiones |
| `settlements` | Registro de liquidaciones y pagos a organizadores |
| `system_parameters` | Parámetros de negocio configurables dinámicamente |

### 4.3 Reutilización de Código

El módulo reutiliza componentes existentes:

- **Middleware de autenticación:** `internal/adapters/http/middleware/auth.go`
- **Token Manager (JWT):** `internal/adapters/redis/token_manager.go`
- **Logger:** `pkg/logger/logger.go`
- **Config:** `pkg/config/config.go`
- **Error handling:** `pkg/errors/errors.go`
- **Repositorios base:** Extiende repositorios existentes (UserRepository, RaffleRepository, PaymentRepository)

---

## 5. Seguridad en Capas

### 5.1 Capa de Presentación (Frontend)

```typescript
// Route protection
<ProtectedRoute requiredRole="super_admin">
  <AdminLayout />
</ProtectedRoute>

// API calls con JWT
const headers = {
  'Authorization': `Bearer ${accessToken}`,
  'Content-Type': 'application/json'
}
```

### 5.2 Capa de Aplicación (HTTP Handlers)

```go
// Middleware chain
admin := router.Group("/api/v1/admin")
admin.Use(authMiddleware.Authenticate())           // Validar JWT
admin.Use(authMiddleware.RequireRole("super_admin")) // Verificar rol
admin.Use(rateLimiter.Limit(10))                    // Rate limiting
admin.Use(auditLogger.LogRequest())                 // Auditoría
{
    admin.GET("/users", userHandler.List)
    // ... más endpoints
}
```

### 5.3 Capa de Dominio (Use Cases)

```go
func (uc *UpdateUserStatusUseCase) Execute(ctx context.Context, req *UpdateStatusRequest) error {
    // Validación de negocio
    if req.UserID == req.AdminID {
        return errors.New("admin cannot suspend themselves")
    }

    // Validación de estado
    if !isValidStatusTransition(user.Status, req.NewStatus) {
        return errors.New("invalid status transition")
    }

    // Lógica de negocio...
}
```

### 5.4 Capa de Persistencia (Database)

```sql
-- Row Level Security (opcional, futura mejora)
CREATE POLICY admin_only ON users
FOR ALL TO sorteos_admin
USING (true);

-- Encriptación de datos sensibles
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Función para encriptar bank account numbers
CREATE FUNCTION encrypt_sensitive(data TEXT) RETURNS TEXT AS $$
BEGIN
    RETURN encode(encrypt(data::bytea, 'encryption_key', 'aes'), 'base64');
END;
$$ LANGUAGE plpgsql;
```

---

## 6. Escalabilidad

### 6.1 Caching Strategy

```
┌─────────────────────────────────────────────────────┐
│ CACHING LAYERS                                      │
├─────────────────────────────────────────────────────┤
│                                                      │
│ 1. Browser Cache (Frontend)                         │
│    - React Query cache (5 min)                      │
│    - Dashboard KPIs cached                          │
│                                                      │
│ 2. CDN Cache (Static Assets)                        │
│    - Nginx cache for /admin bundle.js               │
│                                                      │
│ 3. Application Cache (Redis)                        │
│    - Dashboard KPIs: 1 min TTL                      │
│    - System parameters: 5 min TTL                   │
│    - Company settings: 10 min TTL                   │
│    - Reports: 15 min TTL                            │
│                                                      │
│ 4. Database Cache (PostgreSQL)                      │
│    - Materialized views para reports (refresh 1h)   │
│    - Query result cache                             │
│                                                      │
└─────────────────────────────────────────────────────┘
```

### 6.2 Database Indexing

```sql
-- Índices para queries comunes de admin

-- Búsqueda de usuarios
CREATE INDEX idx_users_search ON users USING gin(to_tsvector('spanish', first_name || ' ' || last_name || ' ' || email));

-- Filtrado de organizadores
CREATE INDEX idx_organizer_profiles_verified ON organizer_profiles(verified, created_at DESC);
CREATE INDEX idx_organizer_profiles_revenue ON organizer_profiles(total_payouts DESC);

-- Settlements por status
CREATE INDEX idx_settlements_status_created ON settlements(status, created_at DESC);

-- Audit logs por fecha y severity
CREATE INDEX idx_audit_logs_severity_date ON audit_logs(severity, created_at DESC);
CREATE INDEX idx_audit_logs_admin_action ON audit_logs(admin_id, action, created_at DESC);
```

### 6.3 Query Optimization

```go
// Paginación eficiente con cursor-based pagination (futura mejora)
type PaginationRequest struct {
    Limit  int
    Cursor string // encoded last_id + last_created_at
}

// Bulk operations
func (r *SettlementRepository) CreateBatch(settlements []*Settlement) error {
    // INSERT multiple rows in single transaction
}
```

---

## 7. Monitoreo y Observabilidad

### 7.1 Métricas Clave

```
┌─────────────────────────────────────────────────────┐
│ MÉTRICAS A MONITOREAR                               │
├─────────────────────────────────────────────────────┤
│                                                      │
│ Performance:                                         │
│ - Dashboard load time (target: <2s)                 │
│ - API response time (p50, p95, p99)                 │
│ - Database query time                               │
│                                                      │
│ Business:                                            │
│ - # de suspensiones de usuarios por día             │
│ - # de settlements aprobados por día                │
│ - # de refunds procesados                           │
│ - # de acciones de admin por tipo                   │
│                                                      │
│ Security:                                            │
│ - Rate limit violations                             │
│ - Failed authentication attempts                    │
│ - Critical audit log events                         │
│                                                      │
│ Errors:                                              │
│ - Error rate por endpoint                           │
│ - Failed payment refunds                            │
│ - Database connection errors                        │
│                                                      │
└─────────────────────────────────────────────────────┘
```

### 7.2 Logging Strategy

```go
// Structured logging con Zap
logger.Info("user suspended",
    zap.Int64("user_id", userID),
    zap.Int64("admin_id", adminID),
    zap.String("reason", reason),
    zap.String("action", "user_suspended"),
    zap.String("severity", "warning"),
)

// Logs de auditoría en DB
auditLog := domain.NewAuditLog("user_suspended").
    WithUser(userID).
    WithAdmin(adminID).
    WithSeverity(domain.AuditSeverityWarning).
    WithEntity("user", userID).
    WithMetadata(map[string]interface{}{
        "status_before": "active",
        "status_after": "suspended",
        "reason": reason,
    }).
    Build()
```

---

## 8. Decisiones Arquitectónicas Clave

### 8.1 ADR-001: Usar Arquitectura Hexagonal

**Decisión:** Mantener arquitectura hexagonal existente para el módulo Almighty.

**Razones:**
- Consistencia con el resto del sistema
- Facilita testing
- Permite cambiar infraestructura sin afectar lógica de negocio

**Consecuencias:**
- Mayor cantidad de código (interfaces, adaptadores)
- Curva de aprendizaje para nuevos developers

### 8.2 ADR-002: Super Admin Único vs Múltiples Roles

**Decisión:** Usar rol único `super_admin` en MVP, preparar para RBAC futuro.

**Razones:**
- Simplicidad en MVP
- Solo el owner necesita acceso inicialmente
- Preparar campo `permissions` para futuro

**Consecuencias:**
- Todos los super_admins tienen mismo nivel de acceso
- Migración futura a RBAC será necesaria

### 8.3 ADR-003: Encriptación de Secrets en DB

**Decisión:** Encriptar credenciales de payment processors y bank info en DB.

**Razones:**
- Seguridad de datos sensibles
- Compliance (PCI DSS para datos de pago)
- Prevenir exposición en backups

**Consecuencias:**
- Performance overhead mínimo
- Complejidad en key management

### 8.4 ADR-004: Audit Log Granular

**Decisión:** Registrar TODAS las acciones de admin en `audit_logs` con metadata JSON.

**Razones:**
- Trazabilidad completa
- Compliance y auditorías
- Debugging de problemas

**Consecuencias:**
- Tabla de audit_logs crece rápidamente
- Necesidad de archivado/purging periódico

### 8.5 ADR-005: Dashboard KPIs con Cache

**Decisión:** Cachear dashboard KPIs en Redis con TTL de 1 minuto.

**Razones:**
- Dashboard es endpoint más consultado
- Cálculos de KPIs son costosos
- 1 minuto de delay es aceptable

**Consecuencias:**
- Dashboard no es 100% real-time
- Complejidad en invalidación de cache

---

## 9. Diagramas Complementarios

### 9.1 Diagrama de Secuencia: Aprobar Settlement

```
Admin          Frontend       API Handler      Use Case        Repository      Database       Email Service
  │               │               │               │               │               │               │
  │─ Click       │               │               │               │               │               │
  │ "Aprobar"    │               │               │               │               │               │
  │──────────────>│               │               │               │               │               │
  │               │─ POST        │               │               │               │               │
  │               │ /settlements/│               │               │               │               │
  │               │ :id/approve  │               │               │               │               │
  │               │──────────────>│               │               │               │               │
  │               │               │─ Authenticate│               │               │               │
  │               │               │ & Authorize  │               │               │               │
  │               │               │──────────────>│               │               │               │
  │               │               │               │─ Get         │               │               │
  │               │               │               │ Settlement   │               │               │
  │               │               │               │──────────────>│               │               │
  │               │               │               │               │─ SELECT      │               │
  │               │               │               │               │──────────────>│               │
  │               │               │               │               │<──────────────│               │
  │               │               │               │<──────────────│               │               │
  │               │               │               │─ Validate    │               │               │
  │               │               │               │ (status=     │               │               │
  │               │               │               │  pending)    │               │               │
  │               │               │               │─ Update      │               │               │
  │               │               │               │ status=      │               │               │
  │               │               │               │ approved     │               │               │
  │               │               │               │──────────────>│               │               │
  │               │               │               │               │─ UPDATE      │               │
  │               │               │               │               │──────────────>│               │
  │               │               │               │               │<──────────────│               │
  │               │               │               │<──────────────│               │               │
  │               │               │               │─ Create      │               │               │
  │               │               │               │ Audit Log    │               │               │
  │               │               │               │──────────────>│               │               │
  │               │               │               │               │─ INSERT      │               │
  │               │               │               │               │──────────────>│               │
  │               │               │               │               │<──────────────│               │
  │               │               │               │<──────────────│               │               │
  │               │               │               │─ Send Email  │               │               │
  │               │               │               │ to Organizer │               │               │
  │               │               │               │───────────────────────────────────────────────>│
  │               │               │               │<───────────────────────────────────────────────│
  │               │               │<──────────────│               │               │               │
  │               │<──────────────│               │               │               │               │
  │<──────────────│               │               │               │               │               │
  │ Toast:        │               │               │               │               │               │
  │ "Settlement   │               │               │               │               │               │
  │  approved!"   │               │               │               │               │               │
```

---

## 10. Referencias y Recursos

### 10.1 Documentación Relacionada

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture - Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)

### 10.2 Tecnologías Utilizadas

- **Go:** https://go.dev/
- **Gin Framework:** https://gin-gonic.com/
- **PostgreSQL:** https://www.postgresql.org/
- **Redis:** https://redis.io/
- **React:** https://react.dev/
- **shadcn/ui:** https://ui.shadcn.com/
- **Tailwind CSS:** https://tailwindcss.com/

---

**Última actualización:** 2025-11-18
**Autor:** Equipo de desarrollo Sorteos.club
