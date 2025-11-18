# STATUS FASE 6 - LIQUIDACIONES (SETTLEMENTS)

**Fecha:** 2025-11-18
**Estado:** ✅ COMPLETADA
**Progreso:** 100% (5/5 use cases)

---

## 📊 RESUMEN EJECUTIVO

Fase 6 completada exitosamente con la implementación de 5 casos de uso para gestión de liquidaciones (settlements) de organizadores. El módulo permite a los administradores:

- **Listar liquidaciones** con filtros avanzados y estadísticas por status
- **Ver detalles completos** de liquidaciones con timeline de eventos
- **Aprobar liquidaciones** con validación de KYC y cuenta bancaria
- **Rechazar liquidaciones** con razón obligatoria
- **Procesar pagos** marcando liquidaciones como pagadas con referencia

**Líneas de código:** ~800 líneas en 5 archivos
**Compilación:** ✅ Sin errores
**Arquitectura:** Hexagonal/Clean Architecture
**Base de datos:** PostgreSQL con GORM

---

## 🎯 CASOS DE USO IMPLEMENTADOS

### 1. **ListSettlementsUseCase** (252 líneas)

**Archivo:** `backend/internal/usecase/admin/settlement/list_settlements.go`

**Funcionalidad:**
- Listar liquidaciones con paginación
- Filtros avanzados: status, organizer_id, raffle_id, date_range, min/max amount, KYC level, search
- JOIN con tablas raffles y users para obtener detalles
- Estadísticas agregadas por status (pending, approved, paid, rejected)
- Cálculo de montos totales por status

**Características clave:**
```go
type ListSettlementsInput struct {
    Page         int
    PageSize     int
    Status       *string // pending, approved, paid, rejected
    OrganizerID  *int64
    RaffleID     *int64
    DateFrom     *string
    DateTo       *string
    MinAmount    *float64
    MaxAmount    *float64
    Search       string // Buscar en raffle title, organizer name
    OrderBy      string
    KYCLevel     *domain.KYCLevel
    PendingOnly  bool
}

type ListSettlementsOutput struct {
    Settlements       []*SettlementWithDetails
    Total             int64
    Page              int
    PageSize          int
    TotalPages        int
    // Estadísticas por status
    TotalPending      int64
    TotalApproved     int64
    TotalPaid         int64
    TotalRejected     int64
    // Montos totales
    TotalPendingAmount  float64
    TotalApprovedAmount float64
    TotalPaidAmount     float64
}
```

**Query complejo con JOIN:**
```go
query := uc.db.Table("settlements").
    Select(`settlements.*,
        raffles.title as raffle_title,
        COALESCE(users.first_name || ' ' || users.last_name, users.email) as organizer_name,
        users.email as organizer_email,
        users.kyc_level as organizer_kyc_level`).
    Joins("LEFT JOIN raffles ON raffles.id = settlements.raffle_id").
    Joins("LEFT JOIN users ON users.id = settlements.organizer_id")
```

**Estadísticas agregadas:**
```go
statsQuery := uc.db.Table("settlements").
    Select("status, COUNT(*) as count, COALESCE(SUM(net_amount), 0) as amount").
    Group("status")
```

---

### 2. **ViewSettlementDetailsUseCase** (218 líneas)

**Archivo:** `backend/internal/usecase/admin/settlement/view_settlement_details.go`

**Funcionalidad:**
- Vista completa 360° de una liquidación específica
- Incluye: settlement, raffle completa, organizador completo
- Resumen de pagos de la rifa (total, succeeded, refunded, revenue)
- Timeline cronológico de eventos (calculated, approved, rejected, paid)
- Información de cuenta bancaria del organizador

**Estructuras de datos:**
```go
type SettlementFullDetails struct {
    Settlement      *SettlementWithDetails
    Raffle          *domain.Raffle
    Organizer       *domain.User
    PaymentsSummary *PaymentsSummary
    Timeline        []*SettlementEvent
    BankAccount     *OrganizerBankAccount
}

type PaymentsSummary struct {
    TotalPayments      int
    SucceededPayments  int
    RefundedPayments   int
    TotalRevenue       float64
    TotalRefunded      float64
    NetRevenue         float64
    PlatformFeePercent float64
    PlatformFeeAmount  float64
}

type SettlementEvent struct {
    Type      string                 `json:"type"`
    Timestamp time.Time              `json:"timestamp"`
    Actor     *string                `json:"actor,omitempty"`
    Details   string                 `json:"details"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
```

**Cálculo de resumen de pagos:**
```go
uc.db.Table("payments").
    Select(`
        COUNT(*) as total_count,
        COUNT(CASE WHEN status = 'succeeded' THEN 1 END) as succeeded_count,
        COUNT(CASE WHEN status = 'refunded' THEN 1 END) as refunded_count,
        COALESCE(SUM(CASE WHEN status = 'succeeded' THEN amount ELSE 0 END), 0) as total_revenue,
        COALESCE(SUM(CASE WHEN status = 'refunded' THEN amount ELSE 0 END), 0) as total_refunded
    `).
    Where("raffle_id = (SELECT uuid FROM raffles WHERE id = ?)", settlement.RaffleID).
    Scan(&paymentStats)
```

**Timeline construction:**
- Evento de cálculo automático
- Evento de aprobación con nombre del admin
- Evento de rechazo con razón
- Evento de pago con referencia y método

---

### 3. **ApproveSettlementUseCase** (166 líneas)

**Archivo:** `backend/internal/usecase/admin/settlement/approve_settlement.go`

**Funcionalidad:**
- Aprobar liquidación pending
- Validación de KYC level (verified o enhanced)
- Validación de cuenta bancaria verificada
- Agregar notas administrativas con timestamp
- Logging crítico de auditoría

**Validaciones de seguridad:**
```go
// Validar estado pending
if settlement.Status != "pending" {
    return nil, errors.New("VALIDATION_FAILED",
        fmt.Sprintf("cannot approve settlement with status %s", settlement.Status), 400, nil)
}

// Validar KYC level
if organizer.KYCLevel != "verified" && organizer.KYCLevel != "enhanced" {
    return nil, errors.New("VALIDATION_FAILED",
        fmt.Sprintf("cannot approve settlement: organizer KYC level is %s, required verified or enhanced", organizer.KYCLevel), 400, nil)
}

// Verificar cuenta bancaria verificada
var bankAccountCount int64
uc.db.Table("organizer_bank_accounts").
    Where("user_id = ? AND verified_at IS NOT NULL", settlement.OrganizerID).
    Count(&bankAccountCount)

if bankAccountCount == 0 {
    return nil, errors.New("VALIDATION_FAILED",
        "cannot approve settlement: organizer has no verified bank account", 400, nil)
}
```

**Actualización con notas:**
```go
timestamp := now.Format("2006-01-02 15:04:05")
newNote := fmt.Sprintf("[%s] Admin ID %d: APPROVED - %s", timestamp, adminID, input.Notes)

updates := map[string]interface{}{
    "status":      "approved",
    "approved_at": now,
    "approved_by": adminID,
    "updated_at":  now,
    "admin_notes": newNote,
}
```

**Logging crítico:**
```go
uc.log.Error("Admin approved settlement",
    logger.Int64("admin_id", adminID),
    logger.Int64("settlement_id", input.SettlementID),
    logger.Int64("organizer_id", settlement.OrganizerID),
    logger.Float64("net_amount", settlement.NetAmount),
    logger.String("action", "admin_approve_settlement"),
    logger.String("severity", "critical"))
```

---

### 4. **RejectSettlementUseCase** (153 líneas)

**Archivo:** `backend/internal/usecase/admin/settlement/reject_settlement.go`

**Funcionalidad:**
- Rechazar liquidación pending o approved
- Razón obligatoria del rechazo
- Notas adicionales opcionales
- Logging crítico de auditoría

**Validaciones:**
```go
// Validar razón obligatoria
if input.Reason == "" {
    return nil, errors.New("VALIDATION_FAILED", "reason is required for rejection", 400, nil)
}

// Validar estado (puede rechazar pending o approved)
if settlement.Status != "pending" && settlement.Status != "approved" {
    return nil, errors.New("VALIDATION_FAILED",
        fmt.Sprintf("cannot reject settlement with status %s", settlement.Status), 400, nil)
}

// No puede rechazar si ya está pagado
if settlement.Status == "paid" {
    return nil, errors.New("VALIDATION_FAILED", "cannot reject paid settlement", 400, nil)
}
```

**Actualización con razón y notas:**
```go
newNote := fmt.Sprintf("[%s] Admin ID %d: REJECTED - Reason: %s", timestamp, adminID, input.Reason)
if input.Notes != "" {
    newNote += fmt.Sprintf(". Notes: %s", input.Notes)
}

updates := map[string]interface{}{
    "status":           "rejected",
    "rejected_at":      now,
    "rejected_by":      adminID,
    "rejection_reason": input.Reason,
    "updated_at":       now,
    "admin_notes":      newNote,
}
```

---

### 5. **ProcessPayoutUseCase** (200 líneas)

**Archivo:** `backend/internal/usecase/admin/settlement/process_payout.go`

**Funcionalidad:**
- Marcar liquidación como pagada (paid)
- Validar estado approved
- Registrar payment reference y method
- Validar métodos de pago permitidos
- Soportar paid amount diferente al net_amount (con warning)

**Validación de métodos de pago:**
```go
validMethods := map[string]bool{
    "wire_transfer":   true,
    "ach":             true,
    "paypal":          true,
    "stripe_connect":  true,
    "manual":          true,
}
if !validMethods[input.PaymentMethod] {
    return nil, errors.New("VALIDATION_FAILED",
        fmt.Sprintf("invalid payment_method: %s", input.PaymentMethod), 400, nil)
}
```

**Validaciones de estado:**
```go
// Debe estar aprobado
if settlement.Status != "approved" {
    return nil, errors.New("VALIDATION_FAILED",
        fmt.Sprintf("cannot process payout for settlement with status %s, must be approved", settlement.Status), 400, nil)
}

// No puede estar ya pagado
if settlement.PaidAt != nil {
    return nil, errors.New("VALIDATION_FAILED", "settlement is already paid", 400, nil)
}

// Verificar cuenta bancaria (seguridad adicional)
var bankAccountCount int64
uc.db.Table("organizer_bank_accounts").
    Where("user_id = ? AND verified_at IS NOT NULL", settlement.OrganizerID).
    Count(&bankAccountCount)

if bankAccountCount == 0 {
    return nil, errors.New("VALIDATION_FAILED",
        "cannot process payout: organizer has no verified bank account", 400, nil)
}
```

**Advertencia si paid amount difiere:**
```go
if paidAmount != settlement.NetAmount {
    uc.log.Error("WARNING: Paid amount differs from net_amount",
        logger.Int64("settlement_id", input.SettlementID),
        logger.Float64("net_amount", settlement.NetAmount),
        logger.Float64("paid_amount", paidAmount),
        logger.String("severity", "warning"))
}
```

**Integración futura con payment providers:**
```go
// TODO: Integrar con payment provider real
// - Si es stripe_connect: hacer transfer a connected account
// - Si es paypal: hacer mass payment
// - Si es wire_transfer/ach: validar con banco (o solo registrar)
var payoutError error

// Simulación de pago
// if input.PaymentMethod == "stripe_connect" {
//     payoutError = uc.stripeService.Transfer(organizerStripeID, paidAmount, input.PaymentReference)
// } else if input.PaymentMethod == "paypal" {
//     payoutError = uc.paypalService.MassPayout(organizerPaypalEmail, paidAmount, input.PaymentReference)
// }
```

---

## 🔧 DETALLES TÉCNICOS

### Arquitectura

**Patrón:** Hexagonal/Clean Architecture
- Use cases en capa de aplicación
- No dependen de frameworks externos
- Reciben dependencias por inyección (db, logger)
- Retornan errores personalizados del paquete `pkg/errors`

**Estructura de archivos:**
```
backend/internal/usecase/admin/settlement/
├── list_settlements.go           (252 líneas)
├── view_settlement_details.go    (218 líneas)
├── approve_settlement.go         (166 líneas)
├── reject_settlement.go          (153 líneas)
└── process_payout.go             (200 líneas)

Total: 5 archivos, ~800 líneas
```

### Base de Datos

**Tabla principal:** `settlements`

Campos utilizados:
- `id` (int64): Primary key
- `raffle_id` (int64): Foreign key a raffles
- `organizer_id` (int64): Foreign key a users
- `total_revenue` (decimal): Revenue total de la rifa
- `platform_fee` (decimal): Comisión de plataforma
- `net_amount` (decimal): Monto neto al organizador
- `status` (varchar): pending, approved, paid, rejected
- `calculated_at` (timestamp): Fecha de cálculo automático
- `approved_at` (timestamp): Fecha de aprobación
- `approved_by` (int64): Admin que aprobó
- `rejected_at` (timestamp): Fecha de rechazo
- `rejected_by` (int64): Admin que rechazó
- `rejection_reason` (text): Razón del rechazo
- `paid_at` (timestamp): Fecha de pago
- `payment_reference` (varchar): Referencia bancaria
- `payment_method` (varchar): Método de pago usado
- `admin_notes` (text): Notas administrativas acumuladas

**Relaciones:**
- `settlements.raffle_id → raffles.id`
- `settlements.organizer_id → users.id`
- `users.id → organizer_bank_accounts.user_id`

**Queries complejos:**

1. **JOIN triple con estadísticas:**
```sql
SELECT settlements.*,
    raffles.title as raffle_title,
    COALESCE(users.first_name || ' ' || users.last_name, users.email) as organizer_name,
    users.email as organizer_email,
    users.kyc_level as organizer_kyc_level
FROM settlements
LEFT JOIN raffles ON raffles.id = settlements.raffle_id
LEFT JOIN users ON users.id = settlements.organizer_id
WHERE settlements.status = 'pending'
ORDER BY settlements.calculated_at DESC
```

2. **Agregación por status:**
```sql
SELECT status,
    COUNT(*) as count,
    COALESCE(SUM(net_amount), 0) as amount
FROM settlements
GROUP BY status
```

3. **Resumen de pagos por rifa:**
```sql
SELECT
    COUNT(*) as total_count,
    COUNT(CASE WHEN status = 'succeeded' THEN 1 END) as succeeded_count,
    COUNT(CASE WHEN status = 'refunded' THEN 1 END) as refunded_count,
    COALESCE(SUM(CASE WHEN status = 'succeeded' THEN amount ELSE 0 END), 0) as total_revenue,
    COALESCE(SUM(CASE WHEN status = 'refunded' THEN amount ELSE 0 END), 0) as total_refunded
FROM payments
WHERE raffle_id = (SELECT uuid FROM raffles WHERE id = ?)
```

### Logging y Auditoría

**Nivel de severidad:** Critical para operaciones financieras

**Eventos auditados:**
- Listar liquidaciones (Info)
- Ver detalles de liquidación (Info)
- Aprobar liquidación (Error/Critical)
- Rechazar liquidación (Error/Critical)
- Procesar pago (Error/Critical)

**Campos en logs:**
```go
logger.Int64("admin_id", adminID)
logger.Int64("settlement_id", settlementID)
logger.Int64("organizer_id", organizerID)
logger.Float64("net_amount", netAmount)
logger.String("action", "admin_approve_settlement")
logger.String("severity", "critical")
```

### Validaciones de Seguridad

1. **Validación de KYC para aprobación:**
   - Organizador debe tener KYC level: verified o enhanced
   - Rechaza si es none, basic, pending

2. **Validación de cuenta bancaria:**
   - Debe existir al menos una cuenta verificada
   - Campo `verified_at` debe ser NOT NULL

3. **Validación de estados:**
   - Solo pending puede ser aprobado
   - Solo pending/approved pueden ser rechazados
   - Solo approved puede ser marcado como paid
   - Paid no puede ser modificado

4. **Validación de métodos de pago:**
   - Lista blanca: wire_transfer, ach, paypal, stripe_connect, manual
   - Rechaza otros métodos

5. **Razón obligatoria para rechazo:**
   - Campo `reason` no puede estar vacío
   - Se registra en `rejection_reason` y `admin_notes`

### TODOs para Integración

**1. Payment Providers:**
```go
// Stripe Connect Transfer
// if input.PaymentMethod == "stripe_connect" {
//     payoutError = uc.stripeService.Transfer(organizerStripeID, paidAmount, input.PaymentReference)
// }

// PayPal Mass Payment
// else if input.PaymentMethod == "paypal" {
//     payoutError = uc.paypalService.MassPayout(organizerPaypalEmail, paidAmount, input.PaymentReference)
// }
```

**2. Email Notifications:**
- Aprobación: Confirmar aprobación, próximos pasos
- Rechazo: Notificar rechazo, razón, acciones correctivas
- Pago: Confirmar pago procesado, comprobante

**3. Platform Fee Configuration:**
```go
platformFeePercent := 10.0 // TODO: Obtener de configuración
```

---

## 🎨 PATRONES DE DISEÑO

### 1. Repository Pattern
- Use cases reciben `*gorm.DB` pero podrían recibir interfaces
- Abstracción sobre acceso a datos

### 2. Command Pattern
- Cada use case es un comando ejecutable
- Input/Output bien definidos
- Execute() como punto de entrada único

### 3. Builder Pattern (implícito)
- Construcción gradual de queries con GORM
- Filtros aplicados condicionalmente

### 4. State Machine (para settlements)
```
pending → approved → paid
  ↓         ↓
rejected  rejected
```

### 5. Audit Trail Pattern
- Timestamped notes en `admin_notes`
- Timeline de eventos
- Logging crítico de operaciones

---

## 📝 ERRORES ENCONTRADOS Y RESUELTOS

### Error 1: Type mismatch en slice de punteros

**Descripción:** `cannot use results (variable of type []SettlementWithDetails) as []*SettlementWithDetails value in struct literal`

**Causa:** `Scan()` retorna slice de structs, pero `ListSettlementsOutput.Settlements` espera slice de punteros.

**Solución:**
```go
// Obtener como slice de valores
var results []SettlementWithDetails
query.Offset(offset).Limit(input.PageSize).Scan(&results)

// Convertir a slice de punteros
settlements := make([]*SettlementWithDetails, len(results))
for i := range results {
    settlements[i] = &results[i]
}
```

**Archivo:** [list_settlements.go:177-181](backend/internal/usecase/admin/settlement/list_settlements.go#L177-L181)

---

## ✅ CRITERIOS DE ACEPTACIÓN

### Funcionales

- [x] Listar liquidaciones con filtros avanzados
- [x] Filtros por status, organizer, raffle, fechas, montos, KYC
- [x] Estadísticas agregadas por status
- [x] Vista completa 360° de liquidación
- [x] Timeline de eventos con actores
- [x] Resumen de pagos de la rifa
- [x] Aprobar con validación de KYC y banco
- [x] Rechazar con razón obligatoria
- [x] Procesar pago con referencia y método
- [x] Validación de métodos de pago permitidos

### Técnicos

- [x] Hexagonal/Clean Architecture
- [x] Compilación sin errores
- [x] Sin imports no utilizados
- [x] Logging con severidad apropiada
- [x] Validaciones de seguridad robustas
- [x] Queries optimizados con JOINs
- [x] TODO markers para integraciones futuras
- [x] Mensajes de error descriptivos

### Seguridad

- [x] Validación de KYC level para aprobaciones
- [x] Validación de cuenta bancaria verificada
- [x] Máquina de estados respetada
- [x] Logging crítico de operaciones financieras
- [x] Razón obligatoria para rechazos
- [x] Payment method validation

---

## 📊 MÉTRICAS DE PROGRESO

### Fase 6
- **Use Cases:** 5/5 (100%)
- **Líneas de código:** ~800
- **Archivos creados:** 5
- **Compilación:** ✅ Exitosa
- **Estado:** ✅ COMPLETADA

### Progreso General Almighty
- **Casos de Uso:** 25/47 (53%)
- **Total Tareas:** 37/185 (20%)
- **Fases Completadas:** 4/8 (Fase 1, 4, 5, 6)
- **Fases Pendientes:** 4 (Fase 2, 3, 7, 8)

---

## 🚀 PRÓXIMOS PASOS

### Fase 7: Análisis y Reportes (7 use cases)

1. **GetPlatformStatisticsUseCase**
   - Métricas globales de plataforma
   - Usuarios, rifas, pagos, revenue

2. **GenerateRevenueReportUseCase**
   - Reporte de ingresos por período
   - Gráficos y tendencias

3. **ListAuditLogsUseCase**
   - Historial completo de acciones admin
   - Filtros por admin, tipo de acción, fecha

4. **ExportDataUseCase**
   - Exportar datos en CSV/Excel
   - Users, raffles, payments, settlements

Y más...

---

## 📚 DOCUMENTACIÓN RELACIONADA

- [ROADMAP_ALMIGHTY.md](ROADMAP_ALMIGHTY.md) - Roadmap completo actualizado
- [STATUS_FASE_5.md](STATUS_FASE_5.md) - Fase anterior (Raffles & Payments)
- [SORTEOS_CONTEXTO_COMPLETO.md](../SORTEOS_CONTEXTO_COMPLETO.md) - Contexto del proyecto

---

**Última actualización:** 2025-11-18
**Responsable:** Claude Code (Almighty Admin Module)
**Estado:** ✅ FASE 6 COMPLETADA - LISTO PARA FASE 7
